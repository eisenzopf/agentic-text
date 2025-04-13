package easy

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// DefaultConfig provides sensible defaults for processor configuration
var DefaultConfig = &Config{
	Provider:    llm.Google,
	Model:       "gemini-2.0-flash",
	MaxTokens:   1024,
	Temperature: 0.2,
}

// Config holds configuration for the easy wrapper
type Config struct {
	// Provider type (Google, OpenAI, etc.)
	Provider llm.ProviderType
	// Model name to use
	Model string
	// MaxTokens limits the response length
	MaxTokens int
	// Temperature controls randomness (0.0-1.0)
	Temperature float64
	// APIKey for the LLM provider (if not set, will use environment variable)
	APIKey string
	// APIKeyEnvVar specifies the environment variable name for the API key
	APIKeyEnvVar string
	// Debug enables debug mode with additional information
	Debug bool
	// Additional provider-specific options
	Options map[string]interface{}
}

// ProcessorWrapper provides a simple interface to use processors
type ProcessorWrapper struct {
	config     *Config
	provider   llm.Provider
	processor  processor.Processor
	procType   string
	procConfig processor.Options
}

// New creates a new processor wrapper with the default configuration
func New(processorType string) (*ProcessorWrapper, error) {
	return NewWithConfig(processorType, DefaultConfig)
}

// NewWithConfig creates a new processor wrapper with the provided configuration
func NewWithConfig(processorType string, config *Config) (*ProcessorWrapper, error) {
	if config == nil {
		config = DefaultConfig
	}

	// Get API key from environment variable if not specified directly
	apiKey := config.APIKey
	if apiKey == "" {
		envVar := config.APIKeyEnvVar
		if envVar == "" {
			// Default environment variable names based on provider
			switch config.Provider {
			case llm.Google:
				envVar = "GEMINI_API_KEY"
			case llm.OpenAI:
				envVar = "OPENAI_API_KEY"
			case llm.Groq:
				envVar = "GROQ_API_KEY"
			case llm.Amazon:
				envVar = "AMAZON_API_KEY"
			default:
				return nil, fmt.Errorf("unknown provider type: %s", config.Provider)
			}
		}

		apiKey = os.Getenv(envVar)
		if apiKey == "" {
			return nil, fmt.Errorf("API key not found in environment variable: %s", envVar)
		}
	}

	// Prepare LLM configuration
	llmConfig := llm.Config{
		APIKey:      apiKey,
		Model:       config.Model,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		Options:     map[string]interface{}{},
	}

	// Copy any additional options
	if config.Options != nil {
		for k, v := range config.Options {
			llmConfig.Options[k] = v
		}
	}

	// Set debug mode if enabled
	if config.Debug {
		llmConfig.Options["debug"] = true
	}

	// Create the provider
	provider, err := llm.NewProvider(config.Provider, llmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Create processor options
	procOptions := processor.Options{
		LLMOptions: llmConfig.Options,
	}

	// Create the processor
	proc, err := processor.Create(processorType, provider, procOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create processor: %w", err)
	}

	return &ProcessorWrapper{
		config:     config,
		provider:   provider,
		processor:  proc,
		procType:   processorType,
		procConfig: procOptions,
	}, nil
}

// Process takes a string input and returns the processed result as a map[string]interface{}
func (w *ProcessorWrapper) Process(input string) (map[string]interface{}, error) {
	if w.processor == nil {
		return nil, errors.New("processor not initialized")
	}

	ctx := context.Background()
	item := data.NewTextProcessItem("input", input, nil)

	result, err := w.processor.Process(ctx, item)
	if err != nil {
		return nil, err
	}

	// Extract processor results based on processor type
	if procInfo, ok := result.ProcessingInfo[w.procType]; ok {
		if resultMap, ok := procInfo.(map[string]interface{}); ok {
			// Clean the response in case it contains JSON in a response field
			return CleanLLMResponse(resultMap), nil
		}
	}

	// Return content as is if no specific processing info is available
	if result.ContentType == "json" {
		if contentMap, ok := result.Content.(map[string]interface{}); ok {
			// Clean the response in case it contains JSON in a response field
			return CleanLLMResponse(contentMap), nil
		}
	}

	return map[string]interface{}{
		"result": result.Content,
	}, nil
}

// ProcessBatch processes multiple inputs in parallel and returns results
func (w *ProcessorWrapper) ProcessBatch(inputs []string, concurrency int) ([]map[string]interface{}, error) {
	if w.processor == nil {
		return nil, errors.New("processor not initialized")
	}

	if concurrency <= 0 {
		concurrency = 2 // Default concurrency
	}

	// Create process items
	items := make([]*data.ProcessItem, len(inputs))
	for i, input := range inputs {
		items[i] = data.NewTextProcessItem(fmt.Sprintf("input-%d", i), input, nil)
	}

	// Create a source
	source := data.NewProcessItemSliceSource(items)

	// Process in parallel
	ctx := context.Background()
	results, err := w.processor.ProcessSource(ctx, source, len(inputs)/concurrency+1, concurrency)
	if err != nil {
		return nil, err
	}

	// Extract results
	outputResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		// Extract processor results based on processor type
		if procInfo, ok := result.ProcessingInfo[w.procType]; ok {
			if resultMap, ok := procInfo.(map[string]interface{}); ok {
				// Clean the response in case it contains JSON in a response field
				outputResults[i] = CleanLLMResponse(resultMap)
				continue
			}
		}

		// Return content as is if no specific processing info is available
		if result.ContentType == "json" {
			if contentMap, ok := result.Content.(map[string]interface{}); ok {
				// Clean the response in case it contains JSON in a response field
				outputResults[i] = CleanLLMResponse(contentMap)
				continue
			}
		}

		outputResults[i] = map[string]interface{}{
			"result": result.Content,
		}
	}

	return outputResults, nil
}

// GetProcessor returns the underlying processor
func (w *ProcessorWrapper) GetProcessor() processor.Processor {
	return w.processor
}

// GetProvider returns the underlying LLM provider
func (w *ProcessorWrapper) GetProvider() llm.Provider {
	return w.provider
}
