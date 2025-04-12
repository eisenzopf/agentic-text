package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// Options holds common configuration for processors
type Options struct {
	// PreProcessOptions holds options for pre-processing
	PreProcessOptions map[string]interface{}
	// LLMOptions holds options for LLM processing
	LLMOptions map[string]interface{}
	// PostProcessOptions holds options for post-processing
	PostProcessOptions map[string]interface{}
}

// TextPreProcessor defines the interface for pre-processing text
type TextPreProcessor interface {
	PreProcess(ctx context.Context, text string) (string, error)
}

// PromptGenerator defines the interface for generating prompts
type PromptGenerator interface {
	GeneratePrompt(ctx context.Context, text string) (string, error)
}

// ResponseHandler defines the interface for handling LLM responses
type ResponseHandler interface {
	HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error)
}

// Processor defines the interface for processors
type Processor interface {
	// GetName returns the name of the processor
	GetName() string

	// GetSupportedContentTypes returns the content types this processor can handle
	GetSupportedContentTypes() []string

	// Process processes a ProcessItem
	Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error)

	// ProcessBatch processes a batch of ProcessItems
	ProcessBatch(ctx context.Context, items []*data.ProcessItem) ([]*data.ProcessItem, error)

	// ProcessSource processes all items from a source
	ProcessSource(ctx context.Context, source data.ProcessItemSource, batchSize, workers int) ([]*data.ProcessItem, error)
}

// BaseProcessor provides a base implementation for processors
type BaseProcessor struct {
	name            string
	contentTypes    []string
	llmClient       llm.Client
	preProcessor    TextPreProcessor
	promptGenerator PromptGenerator
	responseHandler ResponseHandler
	options         Options
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(name string, contentTypes []string, llmClient llm.Client,
	preProcessor TextPreProcessor, promptGenerator PromptGenerator,
	responseHandler ResponseHandler, options Options) *BaseProcessor {

	if contentTypes == nil {
		contentTypes = []string{"text"}
	}

	return &BaseProcessor{
		name:            name,
		contentTypes:    contentTypes,
		llmClient:       llmClient,
		preProcessor:    preProcessor,
		promptGenerator: promptGenerator,
		responseHandler: responseHandler,
		options:         options,
	}
}

// GetName returns the processor name
func (p *BaseProcessor) GetName() string {
	return p.name
}

// GetSupportedContentTypes returns content types this processor can handle
func (p *BaseProcessor) GetSupportedContentTypes() []string {
	return p.contentTypes
}

// Process processes a ProcessItem
func (p *BaseProcessor) Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error) {
	// Validate content type
	contentTypeSupported := false
	for _, ct := range p.contentTypes {
		if ct == item.ContentType {
			contentTypeSupported = true
			break
		}
	}

	if !contentTypeSupported {
		return nil, fmt.Errorf("unsupported content type: %s", item.ContentType)
	}

	// Clone the item to avoid modifying the original
	result, err := item.Clone()
	if err != nil {
		return nil, err
	}

	// Get text content based on the content type
	var textContent string

	if item.ContentType == "text" {
		// Get text content directly
		textContent, err = item.GetTextContent()
		if err != nil {
			return nil, err
		}
	} else if item.ContentType == "json" {
		// For JSON content, either:
		// 1. Use "text" field if available in the JSON
		// 2. Use "response" field if available
		// 3. Or convert the entire JSON to text as fallback
		jsonContent, ok := item.Content.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid JSON content format")
		}

		// Try to extract text from the JSON
		if text, ok := jsonContent["text"].(string); ok {
			textContent = text
		} else if text, ok := jsonContent["response"].(string); ok {
			textContent = text
		} else if originalText, ok := item.Metadata["original_text"].(string); ok {
			// Try to get original text from metadata if available
			textContent = originalText
		} else {
			// Use the first text field we can find
			foundText := false
			for _, value := range jsonContent {
				if text, ok := value.(string); ok {
					textContent = text
					foundText = true
					break
				}
			}

			// If we still don't have text, convert the JSON to string
			if !foundText {
				jsonBytes, err := json.Marshal(jsonContent)
				if err != nil {
					return nil, fmt.Errorf("failed to convert JSON to text: %w", err)
				}
				textContent = string(jsonBytes)
			}
		}
	}

	// Run LLM processing if available
	if p.llmClient != nil {
		// Pre-process if needed
		if p.preProcessor != nil {
			textContent, err = p.preProcessor.PreProcess(ctx, textContent)
			if err != nil {
				return nil, err
			}
		}

		// Generate prompt if needed
		prompt := textContent
		if p.promptGenerator != nil {
			prompt, err = p.promptGenerator.GeneratePrompt(ctx, textContent)
			if err != nil {
				return nil, err
			}
		}

		// Call LLM
		llmResponse, err := p.llmClient.Complete(ctx, prompt, p.options.LLMOptions)
		if err != nil {
			return nil, err
		}

		// Handle response
		if p.responseHandler != nil {
			processedContent, err := p.responseHandler.HandleResponse(ctx, textContent, llmResponse)
			if err != nil {
				return nil, err
			}

			// Update the content with the processed result
			result.Content = processedContent

			// If content is a string, keep content type as text
			// otherwise change to the appropriate type
			if _, ok := processedContent.(string); !ok {
				result.ContentType = "json"
			} else {
				result.ContentType = "text"
			}

			// Add processing info, checking if processor_type already exists in the response
			if contentMap, ok := processedContent.(map[string]interface{}); ok && contentMap["processor_type"] != nil {
				// Use the processor_type from the response
				result.AddProcessingInfo(p.name, processedContent)
			} else {
				// Use the default processor_type
				result.AddProcessingInfo(p.name, map[string]string{
					"processor_type": "base",
				})
			}
		} else {
			// Default behavior: replace content with LLM response
			result.Content = llmResponse

			// If response is a string, assume it's text
			if _, ok := llmResponse.(string); ok {
				result.ContentType = "text"
			} else {
				result.ContentType = "json"
			}

			// Add processing info with default processor_type
			result.AddProcessingInfo(p.name, map[string]string{
				"processor_type": "base",
			})
		}
	} else {
		// Add processing info with default processor_type for non-LLM processing
		result.AddProcessingInfo(p.name, map[string]string{
			"processor_type": "base",
		})
	}

	// Store original text in metadata if not already present
	if _, exists := result.Metadata["original_text"]; !exists {
		if result.Metadata == nil {
			result.Metadata = make(map[string]interface{})
		}
		result.Metadata["original_text"] = textContent
	}

	return result, nil
}

// ProcessBatch processes a batch of items
func (p *BaseProcessor) ProcessBatch(ctx context.Context, items []*data.ProcessItem) ([]*data.ProcessItem, error) {
	results := make([]*data.ProcessItem, len(items))

	for i, item := range items {
		result, err := p.Process(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessSource processes all items from a source
func (p *BaseProcessor) ProcessSource(ctx context.Context, source data.ProcessItemSource, batchSize, workers int) ([]*data.ProcessItem, error) {
	processor := data.NewProcessItemParallelProcessor(source, batchSize, workers)
	defer processor.Close()

	return processor.ProcessAll(ctx, p.Process)
}

// InitializeBuiltInProcessors ensures all built-in processors are registered before use
func InitializeBuiltInProcessors() {
	// This function must be called early in the application to ensure
	// all processor init() functions have run and registered their processors

	// Force loading of processor packages by name
	// The side effect of importing these packages is that their init() functions will run
	// and register themselves with the processor registry

	// We don't need any actual code here, just the import side effects
}

// init runs automatically and calls RegisterBuiltInProcessors to ensure processors are registered
func init() {
	// Make sure our init() function runs after all processors are registered
	// This happens automatically due to Go's package initialization order
}
