# Easy Package

This package provides simplified access to the agentic-text functionality through one-liner functions for common text processing operations.

## Features

- Simple one-liner functions for common text processing tasks
- Default configuration with sensible values
- Support for batch processing and concurrency
- Automatic API key management from environment variables
- Debug mode for troubleshooting

## Usage

### One-line Functions

```go
import (
    "fmt"
    "github.com/eisenzopf/agentic-text/pkg/easy"
)

func main() {
    // Sentiment analysis
    result, err := easy.Sentiment("I absolutely love this product")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Sentiment: %+v\n", result)

    // Intent analysis
    result, err = easy.Intent("I want to cancel my subscription")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Intent: %+v\n", result)

    // Pretty print the result
    prettyJSON, _ := easy.PrettyPrint(result)
    fmt.Println(prettyJSON)
}
```

### Batch Processing

```go
// Define inputs
inputs := []string{
    "I am very disappointed with this service",
    "The product is okay, but nothing special",
    "This is the best experience I have ever had",
}

// Process in parallel (with concurrency of 2)
results, err := easy.ProcessBatchText(inputs, "sentiment", 2)
if err != nil {
    fmt.Println("Batch processing failed:", err)
    return
}

// Process results
for i, result := range results {
    prettyJSON, _ := easy.PrettyPrint(result)
    fmt.Printf("Result for input %d:\n%s\n", i+1, prettyJSON)
}
```

### Custom Configuration

```go
// Define custom configuration
config := &easy.Config{
    Provider:    llm.OpenAI,       // Choose provider: llm.Google, llm.OpenAI, llm.Groq, llm.Amazon
    Model:       "gpt-4",          // Model name varies by provider
    MaxTokens:   512,              // Maximum tokens in response
    Temperature: 0.7,              // Higher for more creative outputs
    Debug:       true,             // Include debug info in results
    Options:     map[string]interface{}{}, // Additional provider-specific options
}

// Use the one-liner with custom config
result, err := easy.ProcessTextWithConfig(
    "I absolutely love this product", 
    "sentiment", 
    config,
)
```

### Available Processors

```go
// List all available processors
processors := easy.ListAvailableProcessors()
fmt.Printf("Available processors: %v\n", processors)
```

Built-in processors include:
- `sentiment`: Analyzes the sentiment of text (positive, negative, neutral)
- `intent`: Identifies the user's intent from the text
- `required_attributes`: Identifies required attributes mentioned in the text
- `get_attributes`: Extracts structured attributes from the text

### Using the ProcessorWrapper Directly

```go
// Create a processor wrapper
wrapper, err := easy.New("sentiment")
if err != nil {
    fmt.Println("Failed to create wrapper:", err)
    return
}

// Process text
result, err := wrapper.Process("I absolutely love this product")
if err != nil {
    fmt.Println("Processing failed:", err)
    return
}

// Batch process multiple inputs
results, err := wrapper.ProcessBatch(inputs, 2) // concurrency of 2
if err != nil {
    fmt.Println("Batch processing failed:", err)
    return
}
``` 