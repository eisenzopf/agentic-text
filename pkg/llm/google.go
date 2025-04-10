package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

// GoogleProvider implements the Provider interface for Google's Vertex AI
type GoogleProvider struct {
	config Config
	client *genai.Client
}

// NewGoogleProvider creates a new Google LLM provider
func NewGoogleProvider(config Config) (*GoogleProvider, error) {
	// Try to get API key from environment variable if not provided
	if config.APIKey == "" {
		config.APIKey = os.Getenv("GEMINI_API_KEY")
		if config.APIKey == "" {
			return nil, errors.New("API key is required for Google provider. Set it in config or GEMINI_API_KEY environment variable")
		}
	}

	if config.Model == "" {
		// Set a default model if none specified
		config.Model = "gemini-1.0-pro"
	}

	// Initialize the Google GenAI client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google GenAI client: %w", err)
	}

	return &GoogleProvider{
		config: config,
		client: client,
	}, nil
}

// Generate implements the Provider interface
func (p *GoogleProvider) Generate(ctx context.Context, prompt string) (string, error) {
	// Call the GenerateContent method with the prompt
	result, err := p.client.Models.GenerateContent(ctx, p.config.Model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("Google API generate error: %w", err)
	}

	// Extract and return the text response
	return result.Text(), nil
}

// GenerateJSON implements the Provider interface
func (p *GoogleProvider) GenerateJSON(ctx context.Context, prompt string, responseStruct interface{}) error {
	// Create a system instruction that tells the model to respond with JSON
	jsonInstruction := &genai.Content{
		Parts: []*genai.Part{
			{Text: "You are a helpful assistant that responds with valid JSON only. No explanations, just JSON."},
		},
		Role: "system",
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: jsonInstruction,
	}

	// Call the GenerateContent method with the JSON instruction
	result, err := p.client.Models.GenerateContent(ctx, p.config.Model, genai.Text(prompt), config)
	if err != nil {
		return fmt.Errorf("Google API JSON generate error: %w", err)
	}

	// Extract the text response and parse it as JSON
	jsonResponse := result.Text()

	// Remove any markdown formatting if present (```json and ```)
	jsonResponse = strings.TrimPrefix(jsonResponse, "```json")
	jsonResponse = strings.TrimPrefix(jsonResponse, "```")
	jsonResponse = strings.TrimSuffix(jsonResponse, "```")
	jsonResponse = strings.TrimSpace(jsonResponse)

	// If debug is enabled, wrap the response with debug info
	if p.config.IsDebugEnabled() {
		// The prompt parameter contains the full interpolated prompt
		if err := WrapWithDebugInfo(ctx, p.config, prompt, jsonResponse, responseStruct); err != nil {
			return err
		}
		return nil
	}

	// Normal behavior (no debug)
	if err := json.Unmarshal([]byte(jsonResponse), responseStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return nil
}

// GetType implements the Provider interface
func (p *GoogleProvider) GetType() ProviderType {
	return Google
}

// GetConfig implements the Provider interface
func (p *GoogleProvider) GetConfig() Config {
	return p.config
}
