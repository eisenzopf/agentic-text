package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// OpenAIProvider implements the Provider interface for OpenAI's API
type OpenAIProvider struct {
	config Config
	// client would typically be the OpenAI API client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config Config) (*OpenAIProvider, error) {
	if config.APIKey == "" {
		return nil, errors.New("API key is required for OpenAI provider")
	}

	if config.Model == "" {
		// Set a default model if none specified
		config.Model = "gpt-4"
	}

	return &OpenAIProvider{
		config: config,
		// Initialize OpenAI API client here
	}, nil
}

// Generate implements the Provider interface
func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	// In a real implementation, this would call the OpenAI API
	// This is a placeholder implementation
	return fmt.Sprintf("OpenAI response to: %s", prompt), nil
}

// GenerateJSON implements the Provider interface
func (p *OpenAIProvider) GenerateJSON(ctx context.Context, prompt string, responseStruct interface{}) error {
	// In a real implementation, this would:
	// 1. Call the OpenAI API with JSON mode enabled
	// 2. Parse the response into the provided struct

	// Placeholder implementation
	_, err := p.Generate(ctx, prompt)
	if err != nil {
		return err
	}

	// Pretend we got valid JSON
	mockJSON := `{"result": "Success", "data": "Sample data from OpenAI"}`

	// If debug is enabled, wrap the response with debug info
	if p.config.IsDebugEnabled() {
		if err := WrapWithDebugInfo(ctx, p.config, prompt, mockJSON, responseStruct); err != nil {
			return err
		}
		return nil
	}

	return json.Unmarshal([]byte(mockJSON), responseStruct)
}

// GetType implements the Provider interface
func (p *OpenAIProvider) GetType() ProviderType {
	return OpenAI
}

// GetConfig implements the Provider interface
func (p *OpenAIProvider) GetConfig() Config {
	return p.config
}
