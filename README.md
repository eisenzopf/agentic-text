# Agentic-Text

A Go library for LLM-powered text processing with pluggable models and data sources.

## Features

- **Simple API**: One-liner functions for common text processing operations through the `easy` package
- **LLM Abstraction**: Support for multiple providers (Google, OpenAI, Groq, Amazon)
- **Data Source Abstraction**: Process text from multiple sources with automatic batching
- **Parallel Processing**: Configurable parallelism and batch size
- **Processor Framework**: Standard interface for text processing operations
- **Extensible Architecture**: Easily add custom processors for specific tasks
- **Standardized Data Containers**: Unified ProcessItem structure for different content types
- **Processing History**: Tracking of processing steps and metadata preservation

## Quick Start

1. Set up your API key environment variable:
   ```bash
   # For Google's Gemini model
   export GEMINI_API_KEY=your_api_key_here
  
   ```

2. Create a new project:
   ```bash
   # Create a project directory
   mkdir agentic-text-example
   cd agentic-text-example
   
   # Initialize a Go module
   go mod init agentic-text-test
   ```

3. Create a main.go file with the following content:
   ```go
   package main

   import (
       "fmt"
       "log"

       "github.com/eisenzopf/agentic-text/pkg/easy"
       // Import the builtin package for processor registration
       _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
   )

   func main() {
       // Simple one-line usage
       result, err := easy.Sentiment("I absolutely love this product")
       if err != nil {
           log.Fatalf("Sentiment analysis failed: %v", err)
       }
       
       // Pretty print the result
       prettyResult, err := easy.PrettyPrint(result)
       if err != nil {
           log.Fatalf("Failed to format result: %v", err)
       }
       fmt.Printf("Sentiment analysis result:\n%s\n\n", prettyResult)
       
       // Try intent detection
       intentResult, err := easy.Intent("I want to cancel my subscription")
       if err != nil {
           log.Fatalf("Intent analysis failed: %v", err)
       }
       
       prettyIntentResult, _ := easy.PrettyPrint(intentResult)
       fmt.Printf("Intent analysis result:\n%s\n\n", prettyIntentResult)
       
       // Batch processing example
       inputs := []string{
           "I am very disappointed with this service",
           "The product is okay, but nothing special",
           "This is the best experience I have ever had",
       }
       
       // Process in parallel with concurrency of 2
       batchResults, err := easy.ProcessBatchText(inputs, "sentiment", 2)
       if err != nil {
           log.Fatalf("Batch processing failed: %v", err)
       }
       
       // Process results
       for i, result := range batchResults {
           prettyJSON, _ := easy.PrettyPrint(result)
           fmt.Printf("Result for input %d: '%s'\n%s\n\n", i+1, inputs[i], prettyJSON)
       }
   }
   ```

4. Install dependencies and run the application:
   ```bash
   # Install dependencies
   go mod tidy
   
   # Run the application
   go run main.go
   ```

## Simple Usage (Recommended)

The `easy` package provides a simplified interface for common operations with sensible defaults. This is the recommended approach for most users.

### One-line Functions

Process text with a single function call:

```go
// Sentiment analysis
result, err := easy.Sentiment("I absolutely love this product")
if err != nil {
    fmt.Println("Sentiment analysis failed:", err)
    return
}

// Intent analysis
result, err := easy.Intent("I want to cancel my subscription immediately")
if err != nil {
    fmt.Println("Intent analysis failed:", err)
    return
}

// Generic processor function
result, err := easy.ProcessText("Extract customer needs from this text", "required_attributes")
if err != nil {
    fmt.Println("Processing failed:", err)
    return
}

// Pretty print the result
jsonStr, _ := easy.PrettyPrint(result)
fmt.Println(jsonStr)
```

### Batch Processing

Process multiple texts in parallel:

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

Configure the processor with your own settings:

```go
package main

import (
    "fmt"
    
    "github.com/eisenzopf/agentic-text/pkg/easy"
    "github.com/eisenzopf/agentic-text/pkg/llm"
)

func main() {
    // Define custom configuration
    config := &easy.Config{
        Provider:    llm.OpenAI,               // Choose provider: llm.Google, llm.OpenAI, llm.Groq, llm.Amazon
        Model:       "gpt-4",                  // Model name varies by provider
        MaxTokens:   512,                      // Maximum tokens in response
        Temperature: 0.7,                      // Higher for more creative outputs
        Debug:       true,                     // Include debug info in results
        Options:     map[string]interface{}{}, // Additional provider-specific options
    }
    
    // Use the one-liner with custom config
    result, err := easy.ProcessTextWithConfig(
        "I absolutely love this product", 
        "sentiment", 
        config,
    )
    if err != nil {
        fmt.Println("Processing failed:", err)
        return
    }
    
    // Pretty print the result
    prettyJSON, _ := easy.PrettyPrint(result)
    fmt.Println(prettyJSON)
}
```

### Available Processors

You can list all available processors:

```go
processors := easy.ListAvailableProcessors()
fmt.Printf("Available processors: %v\n", processors)
```

Built-in processors include:
- `sentiment`: Analyzes the sentiment of text (positive, negative, neutral)
- `intent`: Identifies the user's intent from the text
- `required_attributes`: Identifies required attributes mentioned in the text
- `get_attributes`: Extracts structured attributes from the text
- `keyword_extraction`: Extracts important keywords from text with relevance scores and categories
- `speech_act`: Identifies distinct speech acts within text (questions, requests, statements, etc.)

## Advanced Usage

For more control, you can use the lower-level APIs that power the `easy` package.

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

### Using the Low-Level API

For complete control over the processing pipeline:

```go
package main

import (
    "fmt"
    "context"
    
    "github.com/eisenzopf/agentic-text/pkg/llm"
    "github.com/eisenzopf/agentic-text/pkg/processor"
    "github.com/eisenzopf/agentic-text/pkg/data"
)

func main() {
    // Initialize LLM provider with configuration options
    config := llm.Config{
        APIKey:      "your-api-key", // Use environment variables in production
        Model:       "gemini-2.0-flash",  // Model name varies by provider
        MaxTokens:   1024,           // Maximum tokens in response
        Temperature: 0.2,            // Lower for more deterministic outputs
    }
    
    // Initialize an LLM provider
    provider, err := llm.NewProvider(llm.Google, config)
    if err != nil {
        fmt.Println("Failed to initialize provider:", err)
        return
    }
    
    // Create a processor using the provider
    sentimentProcessor, err := processor.Create("sentiment", provider, processor.Options{})
    if err != nil {
        fmt.Println("Failed to get processor:", err)
        return
    }
    
    // Create a ProcessItem
    item := data.NewTextProcessItem("input-1", "I really enjoyed this product", nil)
    
    // Process the item
    result, err := sentimentProcessor.Process(context.Background(), item)
    if err != nil {
        fmt.Println("Processing failed:", err)
        return
    }
    
    // Get processor data from ProcessingInfo
    var outputData interface{}
    if result.ProcessingInfo != nil {
        for _, procInfo := range result.ProcessingInfo {
            outputData = procInfo
        }
    }
    
    fmt.Println(outputData)
}
```

### Batch Processing with ProcessItems

```go
// Create multiple ProcessItems
items := []*data.ProcessItem{
    data.NewTextProcessItem("input-1", "I am really disappointed with this service", nil),
    data.NewTextProcessItem("input-2", "The product is okay, but nothing special", nil),
    data.NewTextProcessItem("input-3", "This is the best experience I have ever had", nil),
}

// Create a ProcessItemSource from the items
source := data.NewProcessItemSliceSource(items)

// Process all items with parallel processing
// Parameters: context, data source, batch size, concurrency
results, err := sentimentProc.ProcessSource(context.Background(), source, 2, 2)
if err != nil {
    fmt.Println("Batch processing failed:", err)
    return
}

// Process the results
for i, result := range results {
    fmt.Printf("Result for input %d:\n", i+1)
    
    // Access content based on type
    if result.ContentType == "text" {
        text, _ := result.GetTextContent()
        fmt.Printf("Content: %s\n", text)
    }
    
    // Get processor data from ProcessingInfo
    var outputData interface{}
    if result.ProcessingInfo != nil {
        for _, procInfo := range result.ProcessingInfo {
            outputData = procInfo
        }
    }
    
    fmt.Printf("Result data: %v\n", outputData)
}
```

### Creating Custom Processors

You can extend Agentic-Text with custom processors for specialized text processing tasks:

```go
// 1. Define your processor struct
type MyCustomProcessor struct {
    processor.BaseProcessor
}

// 2. Create a constructor function
func NewMyCustomProcessor(provider llm.Provider, options processor.Options) (*MyCustomProcessor, error) {
    p := &MyCustomProcessor{}
    
    // Initialize the base processor
    base, err := processor.NewBaseProcessor("my-custom", provider, options)
    if err != nil {
        return nil, err
    }
    p.BaseProcessor = *base
    
    return p, nil
}

// 3. Implement the Process method
func (p *MyCustomProcessor) Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error) {
    // Validate content type
    if item.ContentType != "text" {
        return nil, fmt.Errorf("processor requires text content type, got %s", item.ContentType)
    }
    
    // Get the text content
    text, err := item.GetTextContent()
    if err != nil {
        return nil, err
    }
    
    // Clone the item to avoid modifying the original
    result, err := item.Clone()
    if err != nil {
        return nil, err
    }
    
    // Generate prompt for the LLM
    prompt := fmt.Sprintf(`Analyze the following text and extract specific information:
Text: %s

Respond with a JSON object containing:
- "value": The primary value extracted from the text
- "score": A confidence score from 0.0 to 1.0

Format your response as valid JSON.`, text)

    // Process with LLM
    var responseData interface{}
    err = p.Provider().GenerateJSON(ctx, prompt, &responseData)
    if err != nil {
        return nil, err
    }
    
    // Store the result in ProcessingInfo
    if data, ok := responseData.(map[string]interface{}); ok {
        result.AddProcessingInfo(p.Name(), data)
    }
    
    return result, nil
}

// 4. Register your processor
func init() {
    processor.Register("my-custom", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
        return NewMyCustomProcessor(provider, options)
    })
}
```

## Pipeline Processing

The `pipeline` package allows you to chain multiple processors together for more complex text analysis workflows:

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/eisenzopf/agentic-text/pkg/llm"
    "github.com/eisenzopf/agentic-text/pkg/processor"
    "github.com/eisenzopf/agentic-text/pkg/pipeline"
    "github.com/eisenzopf/agentic-text/pkg/data"
)

func main() {
    // Initialize an LLM provider
    provider, _ := llm.NewProvider(llm.Google, llm.Config{
        APIKey:      "your-api-key",
        Model:       "gemini-2.0-flash",
        Temperature: 0.2,
    })
    
    // Create a pipeline with multiple processors
    chain, err := pipeline.NewChain(
        []string{"sentiment", "keyword_extraction"},
        provider,
        processor.Options{},
    )
    if err != nil {
        fmt.Println("Failed to create pipeline:", err)
        return
    }
    
    // Process a text item through the pipeline
    item := data.NewTextProcessItem("input-1", "I really enjoyed this product!", nil)
    result, err := chain.Process(context.Background(), item)
    if err != nil {
        fmt.Println("Pipeline processing failed:", err)
        return
    }
    
    // Access results from each processor in the chain
    fmt.Println("Sentiment analysis result:")
    fmt.Println(result.ProcessingInfo["sentiment"])
    
    fmt.Println("Keyword extraction result:")
    fmt.Println(result.ProcessingInfo["keyword_extraction"])
}
```

## Examples

See the [examples](./examples) directory for more detailed examples:

- [Easy Usage](./examples/easy_usage): Demonstrates the simplified `easy` package interface
- [Basic Usage](./examples/basic_usage): Demonstrates basic text processing with different processors
- [ProcessItem Usage](./examples/processitem_usage): Shows how to use the ProcessItem approach for more complex processing
- [Custom Processor](./examples/custom_processor): Explains how to create and use custom processors
- [API Deployment](./examples/api_deployment): Demonstrates deploying processors as a REST API

## Documentation

For more details on specific packages, see their respective README files:
- [pkg/easy/README.md](./pkg/easy/README.md): Simplified interface for common operations
- [pkg/processor/README.md](./pkg/processor/README.md): Core processor framework and interfaces
- [pkg/llm/README.md](./pkg/llm/README.md): LLM provider abstraction
- [pkg/data/README.md](./pkg/data/README.md): Data containers and sources
- [pkg/pipeline/README.md](./pkg/pipeline/README.md): Pipeline processing

## License

MIT 