# LLM Provider Package

This package provides a unified interface for interacting with various Large Language Model providers, abstracting away the differences between their APIs.

## Features

- Consistent interface for multiple LLM providers
- Support for Google (Gemini), OpenAI, Groq, and Amazon Bedrock
- Structured JSON response handling
- Debug mode for capturing prompts and responses
- Configurable parameters for all providers

## Usage

### Creating a Provider

```go
import (
    "github.com/eisenzopf/agentic-text/pkg/llm"
)

// Configure your provider
config := llm.Config{
    APIKey:      "your-api-key", // Or use environment variables
    Model:       "gemini-2.0-flash",
    MaxTokens:   1024,
    Temperature: 0.2,
}

// Create a provider
provider, err := llm.NewProvider(llm.Google, config)
if err != nil {
    // Handle error
}
```

### Generating Text

```go
// Generate text response
response, err := provider.Generate(ctx, "Tell me about Go programming language")
if err != nil {
    // Handle error
}
fmt.Println(response)
```

### Generating Structured Data

```go
// Define a struct for the response
type SentimentResult struct {
    Sentiment string  `json:"sentiment"`
    Score     float64 `json:"score"`
}

// Generate structured JSON response
var result SentimentResult
err := provider.GenerateJSON(ctx, "Analyze the sentiment: I love this product!", &result)
if err != nil {
    // Handle error
}
fmt.Printf("Sentiment: %s, Score: %.2f\n", result.Sentiment, result.Score)
```

### Debug Mode

Enable debug mode to capture prompts and raw responses:

```go
config := llm.Config{
    // ... other settings
    Options: map[string]interface{}{
        "debug": true,
    },
}

// The response will include a "debug" field with prompt and raw response information
```

## Supported Providers

### Google (Gemini)

```go
provider, err := llm.NewGoogleProvider(config)
```

### OpenAI

```go
provider, err := llm.NewOpenAIProvider(config)
```

### Groq

```go
provider, err := llm.NewGroqProvider(config)
```

### Amazon (Bedrock)

```go
provider, err := llm.NewAmazonProvider(config)
```

## Configuration Options

The `Config` struct accepts the following fields:

```go
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
``` 