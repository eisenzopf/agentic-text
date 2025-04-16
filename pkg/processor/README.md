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
	"fmt"
	
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// SentimentResult contains the sentiment analysis results
type SentimentResult struct {
	// Sentiment is the overall sentiment (positive, negative, neutral)
	Sentiment string `json:"sentiment" default:"unknown"`
	// Score is the sentiment score (-1.0 to 1.0)
	Score float64 `json:"score" default:"0.0"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// SentimentPrompt is a prompt generator for sentiment analysis
type SentimentPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *SentimentPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	exampleResult := &SentimentResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert sentiment analysis tool that ONLY outputs valid JSON.

**Objective:** Analyze the sentiment expressed in the provided text.

**Input Text:**
%s

**Instructions:**
1. Determine the primary sentiment: "positive", "negative", or "neutral".
2. Assign a precise sentiment score between -1.0 (most negative) and 1.0 (most positive).
3. Format your entire output as a single, valid JSON object.
4. *** IMPORTANT: Your ENTIRE response must be a single JSON object. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil
}

func init() {
	// Register the sentiment processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"sentiment",        // name
		[]string{"text"},   // contentTypes
		&SentimentResult{}, // resultStruct
		&SentimentPrompt{}, // promptGenerator
		nil,                // no custom initialization needed
		false,              // No struct validation needed by default
	)
}

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

When registering a new processor, you can specify whether to enable struct-level validation:

```go
// Register a processor with validation enabled
processor.RegisterGenericProcessor(
    "my_processor",         // name
    []string{"text", "json"}, // contentTypes
    &MyResultStruct{},       // resultStruct
    &MyPromptGenerator{},    // promptGenerator
    nil,                     // no custom initialization needed
    true,                    // Enable struct-level validation
)
```

### How It Works

The validation system will:

1. Check if the required fields exist in the LLM response
2. Validate that the fields have the expected structure
3. Apply validation before any other transformations

### Custom Validation

You can implement custom validation by adding a method to your result struct with the format:

```go
func (r *MyResultStruct) ValidateFieldName() func(interface{}) interface{} {
    return func(val interface{}) interface{} {
        // Custom validation logic
        // Return the validated value
    }
}
```

Where `FieldName` is the name of the field to validate (with first letter capitalized). 