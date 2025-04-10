# Agentic-Text

A Go library for LLM-powered text processing with pluggable models and data sources.

## Features

- **LLM Abstraction**: Support for multiple providers (Google, Amazon, Groq, etc.)
- **Data Source Abstraction**: Process text from multiple sources with automatic batching
- **Processor Framework**: Standard interface for text processing operations
- **Extensible Architecture**: Easily add custom processors for specific tasks
- **Parallel Processing**: Configurable parallelism and batch size

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/eisenzopf/agentic-text.git
   cd agentic-text
   ```

2. Set up your API key environment variable:
   ```bash
   # For Google's Gemini model
   export GEMINI_API_KEY=your_api_key_here
   
   # For OpenAI
   export OPENAI_API_KEY=your_api_key_here
   ```

3. Run the basic example:
   ```bash
    cd examples/basic_usage
    go run main.go "I absolutely love this product!"
   ```

4. Try different processors:
   ```bash
   # The default processor is "sentiment" if not specified
   go run main.go -processor=sentiment 'I absolutely love this product!'
   
   # Use other processors as they become available in the library
   go run main.go -processor=summarization "Long text to summarize..."
   go run main.go -processor=classification "Text to classify..."
   ```

5. Try batch processing multiple texts:
   ```bash
   go run main.go -batch "I'm really disappointed with this service." "The product is okay." "This is the best experience I've ever had\!"
   ```

6. Override configuration parameters from the command line:
   ```bash
   # Override the provider and model
   go run main.go -provider=openai -model=gpt-3.5-turbo "Analyze this text"
   
   # Change the temperature
   go run main.go -temperature=0.7 "I need more creative results"
   
   # Use a different API key environment variable
   go run main.go -api-key-env=MY_CUSTOM_API_KEY "Test with different credentials"
   
   # Combine multiple overrides
   go run main.go -provider=openai -model=gpt-4 -temperature=0.9 -max-tokens=2048 "Complex analysis needed"
   ```

7. Customize configuration:
   Edit the `config.json` file to change provider, model, or other settings:
   ```json
   {
     "provider": "google",
     "model": "gemini-pro",
     "api_key_env_var": "GOOGLE_API_KEY",
     "max_tokens": 1024,
     "temperature": 0.2
   }
   ```

## Getting Started

### Installation

```bash
go get github.com/eisenzopf/agentic-text
```

### Quick Example

```go
package main

import (
    "fmt"
    "context"
    
    "github.com/eisenzopf/agentic-text/pkg/llm"
    "github.com/eisenzopf/agentic-text/pkg/processor"
)

func main() {
    // Initialize an LLM provider
    provider, err := llm.NewProvider("google")
    if err != nil {
        panic(err)
    }
    
    // Create a processor using the provider
    sentimentProcessor := processor.GetProcessor("sentiment", provider)
    
    // Process a text sample
    result, err := sentimentProcessor.Process(context.Background(), "I really enjoyed this product!")
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
}
```

## Usage

### Configuring an LLM Provider

```go
// Initialize LLM provider with configuration options
config := llm.Config{
    APIKey:      "your-api-key", // Use environment variables in production
    Model:       "gemini-pro",   // Model name varies by provider
    MaxTokens:   1024,           // Maximum tokens in response
    Temperature: 0.2,            // Lower for more deterministic outputs
}

// Create a provider (supported providers: Google, Amazon, Groq, etc.)
provider, err := llm.NewProvider(llm.Google, config)
if err != nil {
    log.Fatalf("Failed to initialize provider: %v", err)
}
```

### Processing Individual Texts

```go
// Create a processor for a specific task
sentimentProcessor, err := processor.GetProcessor("sentiment", provider, processor.Options{})
if err != nil {
    log.Fatalf("Failed to get processor: %v", err)
}

// Process a single text
result, err := sentimentProcessor.Process(context.Background(), "I absolutely love this product!")
if err != nil {
    log.Fatalf("Processing failed: %v", err)
}

// Access the typed result data
sentimentResult, ok := result.Data.(processor.SentimentResult)
if ok {
    fmt.Printf("Sentiment: %s\n", sentimentResult.Sentiment)
    fmt.Printf("Score: %.2f\n", sentimentResult.Score)
}
```

### Batch Processing with Data Sources

```go
// Create a data source from an array of strings
texts := []string{
    "I'm really disappointed with this service.",
    "The product is okay, but nothing special.",
    "This is the best experience I've ever had!",
}
source := data.NewStringsSource(texts)

// Process the entire source with parallel processing
// Parameters: context, data source, batch size, concurrency
results, err := sentimentProcessor.ProcessSource(context.Background(), source, 2, 2)
if err != nil {
    log.Fatalf("Batch processing failed: %v", err)
}

// Process the results
for i, result := range results {
    // Get the original text
    origText := ""
    if item, ok := result.Original.(*data.TextItem); ok {
        origText = item.Content
    }
    
    fmt.Printf("Result for text %d: %v\n", i+1, result.Data)
}
```

## Examples

See the [examples](./examples) directory for more detailed examples:

- [Basic Usage](./examples/basic_usage)
- [Custom Processor](./examples/custom_processor)
- [API Deployment](./examples/api_deployment)

## Documentation

[Full documentation coming soon]

## License

MIT 