# Agentic-Text Documentation

## Overview

Agentic-Text is a Go library for LLM-powered text processing with pluggable models and data sources. It provides a flexible framework for processing text using LLMs with standardized interfaces and built-in support for parallelism and batching.

## Core Components

### LLM Providers

The library abstracts various LLM providers behind a common interface:

```go
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
```

Supported providers:
- Google (Vertex AI)
- Amazon (Bedrock)
- Groq
- OpenAI

### Data Sources

The library provides a flexible way to process text from various sources:

```go
// Source defines the interface for data sources
type Source interface {
    // Next returns the next text item or error when exhausted
    Next(context.Context) (*TextItem, error)
    // Close releases any resources used by the source
    Close() error
}
```

Built-in sources:
- SliceSource (for arrays of TextItems)
- StringsSource (convenience wrapper for string arrays)

### Processors

Processors are the core of the library, handling text analysis tasks:

```go
// Processor defines the interface for text processors
type Processor interface {
    // Process processes a single text item
    Process(ctx context.Context, text string) (*Result, error)
    
    // ProcessItem processes a data.TextItem
    ProcessItem(ctx context.Context, item *data.TextItem) (*Result, error)
    
    // ProcessBatch processes a batch of items
    ProcessBatch(ctx context.Context, items []*data.TextItem) ([]*Result, error)
    
    // ProcessSource processes all items from a source
    ProcessSource(ctx context.Context, source data.Source, batchSize, workers int) ([]*Result, error)
    
    // GetName returns the processor name
    GetName() string
}
```

Each processor follows a standard lifecycle:
1. **Pre-processing**: Prepare the text for the LLM
2. **LLM Interaction**: Send prompts to the LLM
3. **Post-processing**: Process the LLM response

### Pipelines

Pipelines allow chaining multiple processors together:

```go
chain := pipeline.NewChain("my-chain", 
    sentimentProcessor, 
    categoryProcessor,
    summaryProcessor,
)

// Process a text through the entire chain
result, err := chain.Process(ctx, "Your text here")
```

## Performance Features

### Batching

Process multiple texts at once for efficiency:

```go
batchProcessor := data.NewBatchProcessor(source, 10)
results, err := batchProcessor.ProcessAll(ctx, myProcessingFunction)
```

### Parallelism

Process texts concurrently for maximum throughput:

```go
parallelProcessor := data.NewParallelProcessor(source, 10, 4) // batch size 10, 4 workers
results, err := parallelProcessor.ProcessAll(ctx, myProcessingFunction)
```

## Creating Custom Processors

1. Create a struct that embeds BaseProcessor
2. Implement the necessary methods (GeneratePrompt, PostProcess)
3. Register your processor with the registry

Example:

```go
// Create your processor
type MyProcessor struct {
    *processor.BaseProcessor
}

// Implement required methods
func (p *MyProcessor) GeneratePrompt(ctx context.Context, text string) (string, error) {
    return fmt.Sprintf("Process this text: %s", text), nil
}

// Register with the registry
processor.Register("my-processor", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
    base := processor.NewBaseProcessor("my-processor", provider, options)
    return &MyProcessor{BaseProcessor: base}, nil
})
```

## API Development

The library can be easily used to create REST APIs:

```go
http.HandleFunc("/api/process", func(w http.ResponseWriter, r *http.Request) {
    // Get request data
    var req struct {
        Text      string `json:"text"`
        Processor string `json:"processor"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    // Get processor
    proc, _ := processor.GetProcessor(req.Processor, myProvider, processor.Options{})
    
    // Process text
    result, _ := proc.Process(r.Context(), req.Text)
    
    // Return result
    json.NewEncoder(w).Encode(result)
})
```

## Best Practices

1. **Provider Initialization**: Use environment variables for API keys
2. **Error Handling**: Always check errors from LLM calls
3. **Context**: Use Go contexts for timeout handling
4. **Batching**: Use appropriate batch sizes for your use case
5. **Parallelism**: Set workers based on available CPU cores
6. **Processor Design**: Keep processors focused on a single task 