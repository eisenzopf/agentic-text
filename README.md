# Agentic-Text

A Go library for LLM-powered text processing with pluggable models and data sources.

## Features

- **LLM Abstraction**: Support for multiple providers (Google, OpenAI)
- **Data Source Abstraction**: Process text from multiple sources with automatic batching
- **Processor Framework**: Standard interface for text processing operations
- **Extensible Architecture**: Easily add custom processors for specific tasks
- **Parallel Processing**: Configurable parallelism and batch size
- **Standardized Data Containers**: Unified ProcessItem structure for different content types
- **Processing History**: Tracking of processing steps and metadata preservation

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
   go run main.go -processor=sentiment "I absolutely love this product!"
   
   # Use other processors as they become available in the library
   go run main.go -processor=summarization "Long text to summarize..."
   go run main.go -processor=classification "Text to classify..."
   ```

5. Try batch processing multiple texts:
   ```bash
   go run main.go -batch "I'm really disappointed with this service." "The product is okay." "This is the best experience I've ever had!"
   ```

6. Try the ProcessItem-based approach:
   ```bash
   cd examples/processitem_usage
   go run main.go "This is a test of the ProcessItem approach"
   
   # Chain multiple processors together
   go run main.go -secondary=intent "I'm very disappointed with your service. I want to cancel my subscription immediately."
   ```

7. Override configuration parameters from the command line:
   ```bash
   # Override the provider and model
   go run main.go -provider=openai -model=gpt-3.5-turbo "Analyze this text"
   
   # Change the temperature
   go run main.go -temperature=0.7 "I need more creative results"
   
   # Use a different API key environment variable
   go run main.go -api-key-env=MY_CUSTOM_API_KEY "Test with different credentials"
   
   # Enable verbose mode to see LLM prompt and responses
   go run main.go -verbose "Show me the LLM prompt and response"
   
   # Combine multiple overrides
   go run main.go -provider=openai -model=gpt-4 -temperature=0.9 -max-tokens=2048 "Complex analysis needed"
   ```

8. Customize configuration:
   Edit the `config.json` file to change provider, model, or other settings:
   ```json
   {
     "provider": "google",
     "model": "gemini-2.0-flash",
     "api_key_env_var": "GEMINI_API_KEY",
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
    "log"
    
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
        log.Fatalf("Failed to initialize provider: %v", err)
    }
    
    // Create a processor using the provider
    sentimentProcessor, err := processor.Create("sentiment", provider, processor.Options{})
    if err != nil {
        log.Fatalf("Failed to get processor: %v", err)
    }
    
    // Create a ProcessItem
    item := data.NewTextProcessItem("input-1", "I really enjoyed this product!", nil)
    
    // Process the item
    result, err := sentimentProcessor.Process(context.Background(), item)
    if err != nil {
        log.Fatalf("Processing failed: %v", err)
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

## Usage

### Configuring an LLM Provider

```go
// Initialize LLM provider with configuration options
config := llm.Config{
    APIKey:      "your-api-key", // Use environment variables in production
    Model:       "gemini-2.0-flash",  // Model name varies by provider
    MaxTokens:   1024,           // Maximum tokens in response
    Temperature: 0.2,            // Lower for more deterministic outputs
    Options:     map[string]interface{}{"debug": true}, // Optional debug mode
}

// Create a provider (currently supported providers: Google, OpenAI)
provider, err := llm.NewProvider(llm.Google, config)
if err != nil {
    log.Fatalf("Failed to initialize provider: %v", err)
}
```

### Using the ProcessItem Approach

The ProcessItem approach provides a standardized container for data flowing through processors:

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/eisenzopf/agentic-text/pkg/data"
    "github.com/eisenzopf/agentic-text/pkg/processor"
    "github.com/eisenzopf/agentic-text/pkg/pipeline"
)

// Create processors
sentimentProc, err := processor.Create("sentiment", provider, processor.Options{})
if err != nil {
    log.Fatalf("Failed to get processor: %v", err)
}

// Create a ProcessItem with metadata
item := data.NewTextProcessItem("input-1", "I absolutely love this product!", map[string]interface{}{
    "source": "customer-review",
    "timestamp": "2023-07-15T10:30:00Z",
})

// Process the item
result, err := sentimentProc.Process(context.Background(), item)
if err != nil {
    log.Fatalf("Processing failed: %v", err)
}

// Access the results
fmt.Printf("Content Type: %s\n", result.ContentType)

// Access processing history
for procName, procInfo := range result.ProcessingInfo {
    fmt.Printf("Processor: %s, Info: %v\n", procName, procInfo)
}

// Access metadata
fmt.Printf("Metadata: %v\n", result.Metadata)
```

### Batch Processing with ProcessItems

```go
// Create multiple ProcessItems
items := []*data.ProcessItem{
    data.NewTextProcessItem("input-1", "I'm really disappointed with this service.", nil),
    data.NewTextProcessItem("input-2", "The product is okay, but nothing special.", nil),
    data.NewTextProcessItem("input-3", "This is the best experience I've ever had!", nil),
}

// Create a ProcessItemSource from the items
source := data.NewProcessItemSliceSource(items)

// Process all items with parallel processing
// Parameters: context, data source, batch size, concurrency
results, err := sentimentProc.ProcessSource(context.Background(), source, 2, 2)
if err != nil {
    log.Fatalf("Batch processing failed: %v", err)
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

## Benefits of the ProcessItem Approach

The ProcessItem-based approach offers several advantages:

1. **Content Type Flexibility**: Processors can handle different content types (text, JSON, etc.) with type safety
2. **Processing History**: The entire processing history is preserved in the ProcessingInfo field
3. **Metadata Preservation**: Metadata is maintained throughout the processing pipeline
4. **Standardization**: Consistent interface for all processors, regardless of input/output types
5. **Extensibility**: New content types can be added without changing the interface
6. **Debugging Support**: Debug information can be included with the processor results

## Examples

See the [examples](./examples) directory for more detailed examples:

- [Basic Usage](./examples/basic_usage): Demonstrates basic text processing with different processors
- [ProcessItem Usage](./examples/processitem_usage): Shows how to use the ProcessItem approach for more complex processing
- [Custom Processor](./examples/custom_processor): Explains how to create and use custom processors
- [API Deployment](./examples/api_deployment): Demonstrates deploying processors as a REST API

## Documentation

[Full documentation coming soon]

## License

MIT 