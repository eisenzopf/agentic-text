# Processor Package

This package provides a framework for creating text processing pipelines using LLMs.

## Package Structure

The processor package is split into multiple files to make it more maintainable:

- `interfaces.go`: Core interfaces for processors
- `base_processor.go`: Base implementation of processor interface
- `generic_processor.go`: Generic processor with standard response handling
- `response_handler.go`: LLM response handling functionality
- `json_utils.go`: JSON utilities for handling structured data
- `validation.go`: Validation functions for LLM responses
- `registry.go`: Processor registration and creation
- `utils.go`: Common utility functions
- `processor.go`: Initialization and registration logic

## Creating a Custom Processor

Here's a simple example of creating a custom processor:

```go
package myprocessors

import (
	"context"
	
	"github.com/eisenzopf/agentic-text/pkg/processor"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// SentimentResult defines the structure of sentiment analysis results
type SentimentResult struct {
	ProcessorType string  `json:"processor_type"`
	Sentiment     string  `json:"sentiment"`
	Score         float64 `json:"score"`
}

// SentimentPromptGenerator generates prompts for sentiment analysis
type SentimentPromptGenerator struct{}

func (g *SentimentPromptGenerator) GeneratePrompt(ctx context.Context, text string) (string, error) {
	return `Analyze the sentiment of the following text. 
	Return a JSON object with "sentiment" (either "positive", "negative", or "neutral") 
	and "score" (a value between -1.0 and 1.0, where -1.0 is very negative and 1.0 is very positive):
	
	TEXT: ` + text, nil
}

func init() {
	// Register the sentiment processor
	processor.RegisterGenericProcessor(
		"sentiment",
		[]string{"text"},
		&SentimentResult{},
		&SentimentPromptGenerator{},
		nil, // No custom init
	)
}
```

## Using Processors

```go
import (
	"context"
	
	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
	_ "myapp/myprocessors" // Import for side effects (registration)
)

func main() {
	// Create a processor
	p, err := processor.Create("sentiment", llm.NewOpenAIProvider(), processor.Options{})
	if err != nil {
		panic(err)
	}
	
	// Create a process item
	item := &data.ProcessItem{
		ID:          "1",
		Content:     "I really enjoyed this product! It works great!",
		ContentType: "text",
	}
	
	// Process the item
	result, err := p.Process(context.Background(), item)
	if err != nil {
		panic(err)
	}
	
	// Use the result
	// The content will be a SentimentResult struct
	fmt.Printf("Sentiment: %+v\n", result.Content)
}
```

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