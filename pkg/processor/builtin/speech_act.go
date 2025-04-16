package builtin

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// ExampleResult demonstrates the simplified processor pattern
type ExampleResult struct {
	// Category of the text (informational, question, request, etc.)
	Category string `json:"category" default:"unknown"`
	// Complexity score from 1-10
	Complexity float64 `json:"complexity" default:"1.0"`
	// Keywords extracted from the text
	Keywords []string `json:"keywords,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// DefaultValues returns the default values for this result type
// This method can now be very simple using the DefaultsFromStruct helper
func (r *ExampleResult) DefaultValues() map[string]interface{} {
	return processor.DefaultsFromStruct(r)
}

// ExamplePrompt is a prompt generator for the example processor
type ExamplePrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *ExamplePrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	exampleResult := &ExampleResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert at categorizing text.

**Objective:** Analyze the provided text and categorize it according to the structure below.

**Input Text:**
%s

**Instructions:**
1. Determine the category of the text (informational, question, request, etc.)
2. Rate the complexity of the text on a scale from 1.0-10.0
3. Extract up to 5 keywords from the text
4. Format your response as a valid JSON object matching the structure provided.
5. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil
}

// Register the processor with the registry
func init() {
	// Register the example processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"speech_act",     // name
		[]string{"text"}, // contentTypes
		&ExampleResult{}, // resultStruct - default values come from struct tags
		&ExamplePrompt{}, // promptGenerator
		nil,              // no custom initialization needed
		false,            // No struct validation needed by default
	)
}
