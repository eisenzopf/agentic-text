package llm

import (
	"context"
)

// Client defines a simplified interface for interacting with LLM services
type Client interface {
	// Complete sends a prompt to the LLM and returns the response
	Complete(ctx context.Context, prompt string, options map[string]interface{}) (interface{}, error)
}

// ProviderClient implements Client using a Provider
type ProviderClient struct {
	provider Provider
}

// NewProviderClient creates a new client from a Provider
func NewProviderClient(provider Provider) *ProviderClient {
	return &ProviderClient{
		provider: provider,
	}
}

// Complete implements the Client interface
func (c *ProviderClient) Complete(ctx context.Context, prompt string, options map[string]interface{}) (interface{}, error) {
	// If options specify JSON output
	if jsonOutput, ok := options["json_output"].(bool); ok && jsonOutput {
		var responseData interface{}
		err := c.provider.GenerateJSON(ctx, prompt, &responseData)
		return responseData, err
	}

	// Default to text output
	response, err := c.provider.Generate(ctx, prompt)
	return response, err
}
