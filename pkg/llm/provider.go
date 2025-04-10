package llm

import (
	"context"
	"fmt"
)

// ProviderType represents the type of LLM provider
type ProviderType string

const (
	// Google provider type
	Google ProviderType = "google"
	// Amazon provider type
	Amazon ProviderType = "amazon"
	// Groq provider type
	Groq ProviderType = "groq"
	// OpenAI provider type
	OpenAI ProviderType = "openai"
)

// Config holds common configuration for all providers
type Config struct {
	// APIKey for the LLM provider
	APIKey string
	// Model name/ID to use
	Model string
	// MaxTokens limits the response length
	MaxTokens int
	// Temperature controls randomness (0.0-1.0)
	Temperature float64
	// Additional provider-specific options
	Options map[string]interface{}
}

// Provider defines the interface for interacting with LLM providers
type Provider interface {
	// Generate prompts the LLM and returns the generated text
	Generate(ctx context.Context, prompt string) (string, error)

	// GenerateJSON prompts the LLM and returns structured JSON
	GenerateJSON(ctx context.Context, prompt string, responseStruct interface{}) error

	// GetType returns the provider type
	GetType() ProviderType

	// GetConfig returns the provider configuration
	GetConfig() Config
}

// NewProvider creates a new LLM provider based on the type
func NewProvider(providerType ProviderType, config Config) (Provider, error) {
	switch providerType {
	case Google:
		return NewGoogleProvider(config)
	case Amazon:
		return NewAmazonProvider(config)
	case Groq:
		return NewGroqProvider(config)
	case OpenAI:
		return NewOpenAIProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}
