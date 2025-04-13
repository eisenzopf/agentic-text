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
	Attributes []AttributeDefinition `json:"attributes,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// DefaultValues returns the default values for this result type
func (r *RequiredAttributesResult) DefaultValues() map[string]interface{} {
	defaultAttr := []AttributeDefinition{
		{
			FieldName:   "unknown",
			Title:       "Unknown",
			Description: "Unable to determine required attributes from the response",
			Rationale:   "The response did not contain valid attribute definitions",
		},
	}

	return map[string]interface{}{
		"attributes": defaultAttr,
	}
}

// ValidateAttributes returns a transform function for validating attributes
func (r *RequiredAttributesResult) ValidateAttributes() func(interface{}) interface{} {
	defaultAttr := []AttributeDefinition{
		{
			FieldName:   "unknown",
			Title:       "Unknown",
			Description: "Unable to determine required attributes from the response",
			Rationale:   "The response did not contain valid attribute definitions",
		},
	}

	return func(val interface{}) interface{} {
		// Try to convert the value to a slice of attributes
		attributesRaw, ok := val.([]interface{})
		if !ok {
			return defaultAttr
		}

		// If no attributes, return default
		if len(attributesRaw) == 0 {
			return defaultAttr
		}

		// Process each attribute to ensure it has the right structure
		validAttributes := make([]interface{}, 0, len(attributesRaw))
		for _, attrRaw := range attributesRaw {
			attrMap, ok := attrRaw.(map[string]interface{})
			if !ok {
				continue // Skip invalid entries
			}

			// Ensure required fields exist and have values
			fieldName := GetStringValue(attrMap, "field_name")
			if fieldName == "" {
				continue // Skip attributes without a field name
			}

			// Add the validated attribute
			validAttributes = append(validAttributes, attrMap)
		}

		// If no valid attributes were found, use default
		if len(validAttributes) == 0 {
			return defaultAttr
		}

		return validAttributes
	}
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
