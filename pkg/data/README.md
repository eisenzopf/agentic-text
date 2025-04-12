# Standardized Data Processing Interface

This package provides a standardized data interface for processing text and other content types within the system. The goal is to create a unified data container that can handle different types of content while maintaining a consistent API.

## Key Components

### ProcessItem

The `ProcessItem` struct serves as a standardized container for data flowing through processors. It supports:

- Different content types (`text`, `json`, etc.)
- Metadata for contextual information
- Processing history tracking
- Type-safe content access

```go
type ProcessItem struct {
    // ID for tracking the item
    ID string `json:"id"`
    
    // Content holds the actual data (could be string, object, etc.)
    Content interface{} `json:"content"`
    
    // ContentType indicates how to interpret the content
    ContentType string `json:"content_type"`
    
    // Metadata for additional information
    Metadata map[string]interface{} `json:"metadata,omitempty"`
    
    // ProcessingInfo contains history and context of processing operations
    ProcessingInfo map[string]interface{} `json:"processing_info,omitempty"`
}
```

### Source Interface

A single interface for data sources:

```go
// ProcessItemSource defines an interface for sources that provide ProcessItems
type ProcessItemSource interface {
    // NextProcessItem returns the next ProcessItem or error when exhausted
    NextProcessItem(context.Context) (*ProcessItem, error)
    // Close releases any resources used by the source
    Close() error
}
```

### Batch and Parallel Processing

Efficient batch and parallel processors for ProcessItems:

- `ProcessItemBatchProcessor` - For batched processing
- `ProcessItemParallelProcessor` - For parallel multi-thread processing

## Usage Example

Basic usage:

```go
// Create a source
source := data.NewTextStringsProcessItemSource([]string{"Text 1", "Text 2"})

// Process through a processor
results, err := processor.ProcessSource(ctx, source, 10, 2)
```

Direct ProcessItem approach:

```go
// Create a ProcessItem
item := data.NewTextProcessItem("1", "Sample text", nil)

// Process it
result, err := processor.Process(ctx, item)

// Check result content type
if result.ContentType == "json" {
    jsonData := result.Content.(map[string]interface{})
    // Work with structured data
} else {
    text, _ := result.GetTextContent()
    // Work with text
}
```

## Key Benefits

1. **Consistency**: All data flows through the system in a standardized container
2. **Flexibility**: Support for different content types
3. **Type Safety**: Safe content extraction with error handling
4. **Processing History**: Track changes through the processing pipeline
5. **Metadata Preservation**: Context is maintained throughout processing

## Future Enhancements

1. More content type handlers
2. Stronger typing using generics
3. More conversion utilities between common formats 