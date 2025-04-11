package processor

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// AttributeDefinition represents a data attribute definition
type AttributeDefinition struct {
	FieldName   string `json:"field_name"`  // Database field name in snake_case
	Title       string `json:"title"`       // Human readable title
	Description string `json:"description"` // Detailed description of the attribute
	Rationale   string `json:"rationale"`   // Why this attribute is needed
}

// RequiredAttributesProcessor generates required attributes based on questions
type RequiredAttributesProcessor struct {
	// Embed BaseProcessor to inherit all methods
	BaseProcessor
}

// NewRequiredAttributesProcessor creates a new RequiredAttributesProcessor
func NewRequiredAttributesProcessor(provider llm.Provider, options Options) (*RequiredAttributesProcessor, error) {
	p := &RequiredAttributesProcessor{}

	// Create and embed base processor
	base := NewBaseProcessor("required_attributes", provider, options, nil, p, p)
	p.BaseProcessor = *base

	return p, nil
}

// GeneratePrompt implements PromptGenerator interface - generates the attribute analysis prompt
func (p *RequiredAttributesProcessor) GeneratePrompt(ctx context.Context, text string) (string, error) {
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

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *RequiredAttributesProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (*Result, error) {
	// Convert the response data to map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format")
	}

	// Check if debug info exists and preserve it
	var debugInfo interface{}
	if debug, exists := data["debug"]; exists {
		debugInfo = debug
	}

	// Extract attributes
	attributesRaw, ok := data["attributes"].([]interface{})
	if !ok {
		// Create a default attribute if none are found
		attributesRaw = []interface{}{
			map[string]interface{}{
				"field_name":  "unknown",
				"title":       "Unknown",
				"description": "Unable to determine required attributes from the response",
				"rationale":   "The response did not contain valid attribute definitions",
			},
		}
	}

	// Convert to strongly typed attributes
	attributes := make([]AttributeDefinition, 0, len(attributesRaw))
	for _, attrRaw := range attributesRaw {
		attrMap, ok := attrRaw.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		attr := AttributeDefinition{
			FieldName:   attrMap["field_name"].(string),
			Title:       attrMap["title"].(string),
			Description: attrMap["description"].(string),
			Rationale:   attrMap["rationale"].(string),
		}

		// Only add if field_name is valid
		if attr.FieldName != "" {
			attributes = append(attributes, attr)
		}
	}

	// Create result map with attributes
	resultMap := map[string]interface{}{
		"attributes": attributes,
	}

	// Add debug info back if it existed
	if debugInfo != nil {
		resultMap["debug"] = debugInfo
	}

	return &Result{
		Original:  text,
		Processed: text,
		Data:      resultMap,
	}, nil
}

// Register the processor with the registry
func init() {
	Register("required_attributes", func(provider llm.Provider, options Options) (Processor, error) {
		return NewRequiredAttributesProcessor(provider, options)
	})
}
