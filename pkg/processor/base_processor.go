package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

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
		// Check if debug is enabled in options
		debugEnabled := false
		if p.options.LLMOptions != nil {
			if debug, ok := p.options.LLMOptions["debug"].(bool); ok {
				debugEnabled = debug
			}
		}

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

		// Print debug information if enabled
		if debugEnabled {
			DebugLLMInteraction(prompt, "") // Print the prompt before calling LLM
		}

		// Call LLM
		llmResponse, err := p.llmClient.Complete(ctx, prompt, p.options.LLMOptions)
		if err != nil {
			return nil, err
		}

		// Print debug information if enabled
		if debugEnabled {
			DebugLLMInteraction(prompt, llmResponse) // Print full interaction
		}

		// Store debug info in a map if debug is enabled
		var debugInfo map[string]interface{}
		if debugEnabled {
			debugInfo = map[string]interface{}{
				"prompt":       prompt,
				"raw_response": llmResponse,
			}
		}

		// Handle response
		if p.responseHandler != nil {
			processedContent, err := p.responseHandler.HandleResponse(ctx, textContent, llmResponse)
			if err != nil {
				return nil, err
			}

			// Add debug info to processed content if available
			if debugEnabled && debugInfo != nil {
				// If the result is a map, add debug info directly
				if contentMap, ok := processedContent.(map[string]interface{}); ok {
					contentMap["debug"] = debugInfo
					processedContent = contentMap
				} else {
					// For struct responses, we'll handle debug in a different way below
				}
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
				// For struct responses, convert to map first
				// This handles cases like SentimentResult, IntentResult, etc.
				if reflect.TypeOf(processedContent) != nil && reflect.TypeOf(processedContent).Kind() == reflect.Ptr {
					// Use reflection to convert struct to map
					val := reflect.ValueOf(processedContent).Elem()
					if val.Kind() == reflect.Struct {
						structMap := make(map[string]interface{})
						structType := val.Type()

						// First see if struct has a ProcessorType field
						var hasProcessorType bool
						var processorTypeValue string

						// Check each field in the struct
						for i := 0; i < val.NumField(); i++ {
							field := structType.Field(i)

							// Get the field's JSON tag
							tag := field.Tag.Get("json")
							if tag == "" {
								tag = strings.ToLower(field.Name)
							} else {
								tag = strings.Split(tag, ",")[0]
							}

							// Skip if the tag is "-" (meaning don't include in JSON)
							if tag == "-" {
								continue
							}

							// Get the field value
							fieldValue := val.Field(i).Interface()
							structMap[tag] = fieldValue

							// Check if this is the processor_type field
							if tag == "processor_type" {
								hasProcessorType = true
								if strValue, ok := fieldValue.(string); ok {
									processorTypeValue = strValue
								}
							}
						}

						// Add debug info to the struct map if enabled
						if debugEnabled && debugInfo != nil {
							structMap["debug"] = debugInfo
						}

						// If the struct has a processor_type, use it
						if hasProcessorType && processorTypeValue != "" {
							result.AddProcessingInfo(p.name, structMap)
							result.Content = processedContent // Keep the original content
						} else {
							// Add the processor type to the map
							structMap["processor_type"] = p.name
							result.AddProcessingInfo(p.name, structMap)
						}

						return result, nil
					}
				}

				// If not a struct or conversion failed, use the default processor_type
				processingInfo := map[string]interface{}{
					"processor_type": p.name,
				}

				// Add debug info if enabled
				if debugEnabled && debugInfo != nil {
					processingInfo["debug"] = debugInfo
				}

				result.AddProcessingInfo(p.name, processingInfo)
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

			// Add processing info with the proper processor type for non-LLM processing
			processingInfo := map[string]interface{}{
				"processor_type": p.name,
			}

			// Add debug info if enabled
			if debugEnabled && debugInfo != nil {
				processingInfo["debug"] = debugInfo
			}

			result.AddProcessingInfo(p.name, processingInfo)
		}
	} else {
		// Add processing info with the proper processor type for non-LLM processing
		result.AddProcessingInfo(p.name, map[string]interface{}{
			"processor_type": p.name,
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
