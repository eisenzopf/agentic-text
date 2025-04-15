package processor

import (
	"context"

	"github.com/eisenzopf/agentic-text/pkg/data"
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
