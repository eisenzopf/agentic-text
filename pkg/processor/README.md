# Agentic Text Processors

This package contains a framework for creating text processors using LLMs.

## Package Organization

The processor package is organized into two main parts:
- **Core Framework**: The base `processor` package provides interfaces, base implementations, and utilities for building LLM-based text processors
- **Built-in Processors**: The `processor/builtin` package provides ready-to-use implementations for common text analysis tasks

### Using Built-in Processors

To use the built-in processors, import both packages:

```go
import (
    "github.com/eisenzopf/agentic-text/pkg/processor"
    _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin" // Import for side effects (registration)
)

// Now you can create any of the built-in processors
sentimentProc, err := processor.Create("sentiment", provider, options)
```

### Available Built-in Processors

The following processors are included in the `builtin` package:
- **sentiment**: Analyzes the sentiment of text with scores and confidence 
- **intent**: Identifies the primary intent in customer service conversations
- **keyword_extraction**: Extracts important keywords with relevance and categories
- **required_attributes**: Identifies data attributes needed to answer questions
- **get_attributes**: Extracts attribute values from text

## Generic Validation

The package supports a generic validation approach that can be used to validate data returned from LLM responses without needing to implement custom validation methods for each processor.

### Using Generic Validation

When registering a new processor, you can specify validation options to automatically validate specific fields:

```go
// Register a processor with validation for the "attributes" field
processor.RegisterGenericProcessor(
    "my_processor",         // name
    []string{"text", "json"}, // contentTypes
    &MyResultStruct{},       // resultStruct
    &MyPromptGenerator{},    // promptGenerator
    nil,                     // no custom initialization needed
    map[string]interface{}{  // validation options
        "field_name": "attributes",
        "default_value": []MyAttribute{
            {
                ID: "default",
                Value: "Default value used when validation fails",
            },
        },
    },
)
```

### Validation Options

The validation options map supports the following keys:

- `field_name`: The name of the field to validate (required)
- `default_value`: The default value to return if validation fails (optional)

### How It Works

The validation system will:

1. Check if the field exists in the LLM response
2. Validate that the field has the expected structure
3. If validation fails, return the default value
4. Apply the validation before any other transformations

### Custom Validation

You can still implement custom validation by adding a method to your result struct with the format:

```go
func (r *MyResultStruct) ValidateFieldName() func(interface{}) interface{} {
    // Custom validation logic
}
```

Where `FieldName` is the name of the field to validate (with first letter capitalized). 