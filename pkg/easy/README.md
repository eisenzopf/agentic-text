# Easy Package

The `easy` package provides a simplified interface for using Agentic-Text processors. It offers both one-liner functions and a more flexible wrapper approach.

## Quick Start

### One-liner approach

Process text with a single function call:

```go
// Simple sentiment analysis
result, err := easy.Sentiment("I absolutely love this product!")
if err != nil {
    log.Fatalf("Sentiment analysis failed: %v", err)
}

// Or use the generic processor function
result, err := easy.ProcessText("I want to cancel my subscription", "intent")
if err != nil {
    log.Fatalf("Intent analysis failed: %v", err)
}

// Pretty print the result
jsonStr, _ := easy.PrettyPrint(result)
fmt.Println(jsonStr)
```

### Using the wrapper

More control with the wrapper approach:

```go
// Create a processor wrapper
wrapper, err := easy.New("sentiment")
if err != nil {
    log.Fatalf("Failed to create wrapper: %v", err)
}

// Process text
result, err := wrapper.Process("I absolutely love this product!")
if err != nil {
    log.Fatalf("Processing failed: %v", err)
}
```

### Custom configuration

Configure the processor with your own settings:

```go
// Define custom configuration
config := &easy.Config{
    Provider:    llm.Google, // or llm.OpenAI, llm.Groq, llm.Amazon
    Model:       "gemini-2.0-flash",
    MaxTokens:   512,
    Temperature: 0.3,
    Debug:       true, // Include debug info in results
}

// Create wrapper with custom config
wrapper, err := easy.NewWithConfig("sentiment", config)
if err != nil {
    log.Fatalf("Failed to create wrapper: %v", err)
}

// Or use the one-liner with config
result, err := easy.ProcessTextWithConfig("I absolutely love this product!", "sentiment", config)
if err != nil {
    log.Fatalf("Processing failed: %v", err)
}
```

### Batch processing

Process multiple texts in parallel:

```go
// Define inputs
inputs := []string{
    "I'm very disappointed with this service.",
    "The product is okay, but nothing special.",
    "This is the best experience I've ever had!",
}

// Process in parallel (with concurrency of 2)
results, err := easy.ProcessBatchText(inputs, "sentiment", 2)
if err != nil {
    log.Fatalf("Batch processing failed: %v", err)
}

// Process results
for i, result := range results {
    fmt.Printf("Result for input %d: %v\n", i+1, result)
}
```

## Available Processors

You can list all available processors:

```go
processors := easy.ListAvailableProcessors()
fmt.Printf("Available processors: %v\n", processors)
```

## Default Configuration

The package comes with sensible defaults:

```go
var DefaultConfig = &Config{
    Provider:    llm.Google,
    Model:       "gemini-2.0-flash",
    MaxTokens:   1024,
    Temperature: 0.2,
}
```

## Environment Variables

The package automatically looks for API keys in environment variables:

- `GEMINI_API_KEY` for Google's Gemini
- `OPENAI_API_KEY` for OpenAI
- `GROQ_API_KEY` for Groq
- `AMAZON_API_KEY` for Amazon

You can also specify a custom environment variable:

```go
config := &easy.Config{
    Provider:     llm.OpenAI,
    APIKeyEnvVar: "MY_CUSTOM_API_KEY",
}
```

## Examples

For complete examples, see the [examples/easy_usage](../../examples/easy_usage) directory. 