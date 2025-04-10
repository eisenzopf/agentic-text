# Agentic-Text

A Go library for LLM-powered text processing with pluggable models and data sources.

## Features

- **LLM Abstraction**: Support for multiple providers (Google, Amazon, Groq, etc.)
- **Data Source Abstraction**: Process text from multiple sources with automatic batching
- **Processor Framework**: Standard interface for text processing operations
- **Extensible Architecture**: Easily add custom processors for specific tasks
- **Parallel Processing**: Configurable parallelism and batch size

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

## Examples

See the [examples](./examples) directory for more detailed examples:

- [Basic Usage](./examples/basic_usage)
- [Custom Processor](./examples/custom_processor)
- [API Deployment](./examples/api_deployment)

## Documentation

[Full documentation coming soon]

## License

MIT 