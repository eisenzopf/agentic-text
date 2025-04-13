package builtin

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// AttributeResult contains the extracted attributes
type AttributeResult struct {
	// Attributes is an array of extracted attributes
	Attributes []Attribute `json:"attributes,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// DefaultValues returns the default values for this result type
func (r *AttributeResult) DefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"attributes": []Attribute{},
	}
}

// ValidateAttributes returns a transform function for validating attributes
func (r *AttributeResult) ValidateAttributes() func(interface{}) interface{} {
	return func(val interface{}) interface{} {
		// Try to convert to array of attributes
		attrs, ok := val.([]interface{})
		if !ok {
			return []Attribute{}
		}

		// Validate each attribute
		validAttrs := make([]interface{}, 0, len(attrs))
		for _, attr := range attrs {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				// Ensure it has a field_name
				fieldName := processor.GetStringValue(attrMap, "field_name")
				if fieldName != "" {
					validAttrs = append(validAttrs, attrMap)
				}
			}
		}

		return validAttrs
	}
}

// Attribute represents a single extracted attribute
type Attribute struct {
	// FieldName is the name of the attribute
	FieldName string `json:"field_name"`
	// Value is the extracted value
	Value string `json:"value"`
	// Confidence is the confidence level for this specific attribute
	Confidence float64 `json:"confidence"`
	// Explanation provides context for this specific attribute
	Explanation string `json:"explanation"`
}

// AttributePrompt is a prompt generator for attribute extraction
type AttributePrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *AttributePrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	return fmt.Sprintf(`**Role:** You are an expert at extracting structured information from text.

**Objective:** Analyze the provided text and extract relevant attributes and their values.

**Input Text:**
%s

**Instructions:**
1. Carefully read and interpret the Input Text.
2. If the input appears to be JSON containing required attributes, use those as a guide to extract values.
3. Extract any relevant attributes and their values.
4. For each attribute, provide:
   - A field name (in snake_case)
   - The extracted value
   - A confidence score (0.0 to 1.0)
   - A brief explanation
5. Assign an overall confidence score for the extraction.
6. Provide a brief overall explanation of how the attributes were determined.
7. Format your entire output as a single, valid JSON object.
8. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure:**
{
  "attributes": [
    {
      "field_name": "attribute_name",
      "value": "extracted_value",
      "confidence": 0.0,
      "explanation": "..."
    },
    ...
  ]
}`, text), nil
}

// Register the processor with the registry
func init() {
	// Register the attribute processor using the generic processor registration with validation
	processor.RegisterGenericProcessor(
		"get_attributes",         // name
		[]string{"text", "json"}, // contentTypes
		&AttributeResult{},       // resultStruct
		&AttributePrompt{},       // promptGenerator
		nil,                      // no custom initialization needed
		map[string]interface{}{ // validation options
			"field_name":    "attributes",
			"default_value": []Attribute{},
		},
	)
}
