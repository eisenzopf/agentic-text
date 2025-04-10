package llm

import (
	"context"
	"encoding/json"
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

// IsDebugEnabled checks if debug mode is enabled in the config
func (c Config) IsDebugEnabled() bool {
	if c.Options == nil {
		return false
	}
	if debug, ok := c.Options["debug"].(bool); ok {
		return debug
	}
	return false
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

// ExtractJSONResponse attempts to parse a raw response as JSON and extract structured data
// If the response is not valid JSON, it creates a normalized structure with the raw text
func ExtractJSONResponse(rawResponse string) (map[string]interface{}, error) {
	var responseMap map[string]interface{}

	// Try to parse as JSON first
	err := json.Unmarshal([]byte(rawResponse), &responseMap)
	if err != nil {
		// If it's not valid JSON, create a standardized response
		responseMap = map[string]interface{}{
			"response": rawResponse,
		}
		return responseMap, nil
	}

	return responseMap, nil
}

// WrapWithDebugInfo adds debug information to the response data if debug is enabled
// This is a helper function that can be used by all provider implementations
func WrapWithDebugInfo(ctx context.Context, config Config, prompt string, rawResponse string, responseStruct interface{}) error {
	if !config.IsDebugEnabled() {
		return nil
	}

	// Extract or create a standardized response map
	responseMap, err := ExtractJSONResponse(rawResponse)
	if err != nil {
		return fmt.Errorf("failed to extract JSON response: %w", err)
	}

	// Create debug info map with the actual prompt sent to the LLM
	debugInfo := map[string]interface{}{
		"prompt":       prompt,
		"raw_response": rawResponse,
		"model":        config.Model,
	}

	// Add debug info to the response map
	responseMap["debug"] = debugInfo

	// Marshal back to JSON and unmarshal into the original responseStruct
	debugJSON, err := json.Marshal(responseMap)
	if err != nil {
		return fmt.Errorf("failed to marshal debug JSON: %w", err)
	}

	if err := json.Unmarshal(debugJSON, responseStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON with debug info: %w", err)
	}

	return nil
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
