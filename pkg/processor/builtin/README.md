# Built-in Text Processors

This package provides ready-to-use implementations of common text processing tasks using the processor framework.

## Available Processors

The following processors are included:

### `sentiment`

Analyzes the sentiment of text, providing:
- Overall sentiment (positive, negative, neutral)
- Sentiment score (-1.0 to 1.0)
- Confidence level (0.0 to 1.0)
- Key sentiment words from the text

Example usage:
```go
import (
    "github.com/eisenzopf/agentic-text/pkg/processor"
    _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
)

// Create the sentiment processor
sentimentProc, err := processor.Create("sentiment", provider, options)
if err != nil {
    // Handle error
}

// Process text
result, err := sentimentProc.Process(ctx, item)
```

### `intent`

Identifies the primary intent in customer service conversations:
- Label name (human-readable intent label)
- Label (machine-readable version)
- Description (concise explanation of the intent)

### `keyword_extraction`

Extracts important keywords from text with:
- Term (the keyword itself)
- Relevance (0.0 to 1.0)
- Category (e.g., "topic", "person", "location")

### `required_attributes`

Identifies data attributes needed to answer a set of questions:
- Field name (machine-readable name)
- Title (human-readable title)
- Description (what the attribute represents)
- Rationale (why the attribute is needed)

### `get_attributes`

Extracts attribute values from text:
- Field name (attribute name)
- Value (extracted value)
- Confidence (confidence level)
- Explanation (context for the extraction)

## Importing

To use these processors, simply import this package for its side effects:

```go
import _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
```

This will register all builtin processors with the processor registry, making them available through `processor.Create()`. 