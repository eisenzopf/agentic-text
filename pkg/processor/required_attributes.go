package processor

import (
	"context"
	"fmt"
)

// AttributeDefinition represents a data attribute definition
type AttributeDefinition struct {
	FieldName   string `json:"field_name"`  // Database field name in snake_case
	Title       string `json:"title"`       // Human readable title
	Description string `json:"description"` // Detailed description of the attribute
	Rationale   string `json:"rationale"`   // Why this attribute is needed
}

// RequiredAttributesResult contains the required attributes results
type RequiredAttributesResult struct {
	// Attributes is an array of required attributes
	Attributes []AttributeDefinition `json:"attributes"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// RequiredAttributesPrompt is a prompt generator for required attributes
type RequiredAttributesPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *RequiredAttributesPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
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
{
  "attributes": [
    {
      "field_name": "example_field",  // Database field name in snake_case
      "title": "Example Field",       // Human readable title
      "description": "Description of what this field represents",
      "rationale": "Why this field is needed to answer the questions"
    }
  ]
}`, text), nil
}

// Register the processor with the registry
func init() {
	// Register the required attributes processor using the generic processor registration
	RegisterGenericProcessor(
		"required_attributes",       // name
		[]string{"text", "json"},    // contentTypes
		&RequiredAttributesResult{}, // resultStruct
		&RequiredAttributesPrompt{}, // promptGenerator
		nil,                         // no custom initialization needed
	)
}
