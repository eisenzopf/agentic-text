package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// Result represents the result of processing text
type Result struct {
	// Original is the original input text or item
	Original interface{}
	// Processed is the processed text (if any)
	Processed string
	// Data holds structured data extracted from the text
	Data interface{}
	// Error holds any processing error
	Error error
}

// Options holds common configuration for processors
type Options struct {
	// PreProcessOptions holds options for pre-processing
	PreProcessOptions map[string]interface{}
	// LLMOptions holds options for LLM processing
	LLMOptions map[string]interface{}
	// PostProcessOptions holds options for post-processing
	PostProcessOptions map[string]interface{}
}

// Processor defines the interface for text processors
type Processor interface {
	// Process processes a single text item
	Process(ctx context.Context, text string) (*Result, error)

	// ProcessItem processes a data.TextItem
	ProcessItem(ctx context.Context, item *data.TextItem) (*Result, error)

	// ProcessBatch processes a batch of items
	ProcessBatch(ctx context.Context, items []*data.TextItem) ([]*Result, error)

	// ProcessSource processes all items from a source
	ProcessSource(ctx context.Context, source data.Source, batchSize, workers int) ([]*Result, error)

	// GetName returns the processor name
	GetName() string
}

// BaseProcessor provides common functionality for processors
type BaseProcessor struct {
	name     string
	provider llm.Provider
	options  Options
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(name string, provider llm.Provider, options Options) *BaseProcessor {
	return &BaseProcessor{
		name:     name,
		provider: provider,
		options:  options,
	}
}

// GetName returns the processor name
func (p *BaseProcessor) GetName() string {
	return p.name
}

// PreProcess implements common pre-processing logic
func (p *BaseProcessor) PreProcess(_ context.Context, text string) (string, error) {
	// Default implementation returns the text unchanged
	return text, nil
}

// PostProcess implements common post-processing logic
func (p *BaseProcessor) PostProcess(_ context.Context, text string, responseData interface{}) (*Result, error) {
	// Default implementation returns the text and data as-is
	return &Result{
		Processed: text,
		Data:      responseData,
	}, nil
}

// Process processes a single text item
func (p *BaseProcessor) Process(ctx context.Context, text string) (*Result, error) {
	// Pre-process the text
	processedText, err := p.PreProcess(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("pre-processing error: %w", err)
	}

	// Generate the prompt
	prompt, err := p.GeneratePrompt(ctx, processedText)
	if err != nil {
		return nil, fmt.Errorf("prompt generation error: %w", err)
	}

	// Get result from LLM
	var responseData interface{}
	err = p.provider.GenerateJSON(ctx, prompt, &responseData)
	if err != nil {
		return nil, fmt.Errorf("LLM error: %w", err)
	}

	// Post-process the result
	result, err := p.PostProcess(ctx, processedText, responseData)
	if err != nil {
		return nil, fmt.Errorf("post-processing error: %w", err)
	}

	// Set the original text in the result
	result.Original = text

	return result, nil
}

// GeneratePrompt generates the prompt for the LLM
func (p *BaseProcessor) GeneratePrompt(_ context.Context, text string) (string, error) {
	// This should be overridden by concrete implementations
	return fmt.Sprintf("Process the following text: %s", text), nil
}

// ProcessItem processes a data.TextItem
func (p *BaseProcessor) ProcessItem(ctx context.Context, item *data.TextItem) (*Result, error) {
	result, err := p.Process(ctx, item.Content)
	if err != nil {
		return nil, err
	}

	// Replace the generated original with the actual item
	result.Original = item

	return result, nil
}

// ProcessBatch processes a batch of items
func (p *BaseProcessor) ProcessBatch(ctx context.Context, items []*data.TextItem) ([]*Result, error) {
	results := make([]*Result, len(items))

	for i, item := range items {
		result, err := p.ProcessItem(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessSource processes all items from a source
func (p *BaseProcessor) ProcessSource(ctx context.Context, source data.Source, batchSize, workers int) ([]*Result, error) {
	processor := data.NewParallelProcessor(source, batchSize, workers)
	defer processor.Close()

	// Convert data.TextItem processor to Result processor
	itemProcessor := func(ctx context.Context, item *data.TextItem) (*data.TextItem, error) {
		result, err := p.ProcessItem(ctx, item)
		if err != nil {
			return nil, err
		}

		// Pack the result into the metadata
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		if item.Metadata == nil {
			item.Metadata = make(map[string]interface{})
		}
		item.Metadata["result"] = string(resultJSON)

		return item, nil
	}

	// Process all items
	processedItems, err := processor.ProcessAll(ctx, itemProcessor)
	if err != nil {
		return nil, err
	}

	// Extract results from metadata
	results := make([]*Result, len(processedItems))
	for i, item := range processedItems {
		resultJSON, ok := item.Metadata["result"].(string)
		if !ok {
			return nil, fmt.Errorf("missing result in item metadata")
		}

		var result Result
		if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
			return nil, err
		}

		results[i] = &result
	}

	return results, nil
}
