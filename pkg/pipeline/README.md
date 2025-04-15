# Pipeline Package

This package provides tools for chaining multiple text processors together to create complex processing workflows.

## Features

- Chain processors together in sequence
- Process individual items or batches
- Support for data sources and parallel processing
- Error handling and propagation

## Usage

### Creating a Chain

```go
import (
    "github.com/eisenzopf/agentic-text/pkg/llm"
    "github.com/eisenzopf/agentic-text/pkg/processor"
    "github.com/eisenzopf/agentic-text/pkg/pipeline"
)

// Create processors
provider, _ := llm.NewProvider(llm.Google, config)
sentimentProc, _ := processor.Create("sentiment", provider, processor.Options{})
keywordProc, _ := processor.Create("keyword_extraction", provider, processor.Options{})
summaryProc, _ := processor.Create("summarize", provider, processor.Options{})

// Create a chain
chain := pipeline.NewChain(
    "my-analysis-chain",
    sentimentProc,
    keywordProc,
    summaryProc,
)
```

### Processing a Single Item

```go
import (
    "context"
    "fmt"
    
    "github.com/eisenzopf/agentic-text/pkg/data"
)

// Create an item to process
item := data.NewTextProcessItem("1", "I really enjoyed this product! It works great!", nil)

// Process through the chain
ctx := context.Background()
result, err := chain.Process(ctx, item)
if err != nil {
    // Handle error
}

// Access results
fmt.Printf("Processing result: %+v\n", result.Content)
fmt.Printf("Processing info: %+v\n", result.ProcessingInfo)
```

### Processing Multiple Items

```go
// Create multiple items
items := []*data.ProcessItem{
    data.NewTextProcessItem("1", "I really enjoyed this product!", nil),
    data.NewTextProcessItem("2", "This service is terrible.", nil),
    data.NewTextProcessItem("3", "The product is okay, nothing special.", nil),
}

// Process through the chain
results, err := chain.ProcessBatch(ctx, items)
if err != nil {
    // Handle error
}

// Access results
for i, result := range results {
    fmt.Printf("Result %d: %+v\n", i+1, result.ProcessingInfo)
}
```

### Processing from a Data Source

```go
// Create a source
source := data.NewTextStringsProcessItemSource([]string{
    "I really enjoyed this product!",
    "This service is terrible.",
    "The product is okay, nothing special.",
})

// Process the source with parallelism
// Parameters: context, source, batch size, number of workers
results, err := chain.ProcessSource(ctx, source, 10, 2)
if err != nil {
    // Handle error
}

// Access results
for i, result := range results {
    fmt.Printf("Result %d: %+v\n", i+1, result.ProcessingInfo)
}
``` 