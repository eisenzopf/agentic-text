package builtin

import (
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

// Register the processor with the registry
func init() {
	processor.NewBuilder("get_attributes").
		WithStruct(&AttributeResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert at extracting structured information from text").
		WithObjective("Analyze the provided text and extract relevant attributes and their values").
		WithInstructions(
			"Carefully read and interpret the Input Text",
			"If the input appears to be JSON containing required attributes, use those as a guide to extract values",
			"Extract any relevant attributes and their values based on the required structure",
			"For each attribute, provide a field name (in snake_case), the extracted value, a confidence score (0.0 to 1.0), and a brief explanation",
			"Assign an overall confidence score for the extraction",
			"Provide a brief overall explanation of how the attributes were determined",
			"Format your entire output as a single, valid JSON object",
		).
		Register()
}
