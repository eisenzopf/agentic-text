package builtin

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// AttributeDefinition represents a data attribute definition
type AttributeDefinition struct {
	FieldName   string `json:"field_name" default:"unknown"`                                                    // Database field name in snake_case
	Title       string `json:"title" default:"Unknown"`                                                         // Human readable title
	Description string `json:"description" default:"Unable to determine required attributes from the response"` // Detailed description of the attribute
	Rationale   string `json:"rationale" default:"The response did not contain valid attribute definitions"`    // Why this attribute is needed
}

// RequiredAttributesResult contains the required attributes results
type RequiredAttributesResult struct {
	// Attributes is an array of required attributes
	Attributes []AttributeDefinition `json:"attributes,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// RequiredAttributesPrompt is a prompt generator for required attributes
type RequiredAttributesPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *RequiredAttributesPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	exampleResult := &RequiredAttributesResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert data analyst that ONLY outputs valid JSON.

**Objective:** Analyze the provided questions and determine what data attributes would be required to answer them accurately.

**Input Questions:**
%s

**Instructions:**
1. Carefully read and interpret the Input Questions.
2. Identify all data attributes needed to answer these questions.
3. For each attribute, provide:
   - A machine-readable field name in snake_case
   - A human-readable title
   - A clear description of what the attribute represents
   - A rationale explaining why this attribute is needed
4. Format your entire output as a single, valid JSON object conforming to the structure below.
5. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil
}

// init registers the processor
func init() {
	// Register the standard processor using the validation approach
	processor.RegisterGenericProcessor(
		"required_attributes",       // name
		[]string{"text", "json"},    // contentTypes
		&RequiredAttributesResult{}, // resultStruct
		&RequiredAttributesPrompt{}, // promptGenerator
		nil,                         // no custom initialization needed
		false,                       // No struct validation needed by default
	)
}
