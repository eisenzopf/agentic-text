package builtin

import (
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

// Register the processor with the registry
func init() {
	processor.NewBuilder("required_attributes").
		WithStruct(&RequiredAttributesResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert data analyst that ONLY outputs valid JSON").
		WithObjective("Analyze the provided questions and determine what data attributes would be required to answer them accurately").
		WithInstructions(
			"Carefully read and interpret the Input Questions",
			"Identify all data attributes needed to answer these questions",
			"For each attribute, provide a machine-readable field name in snake_case",
			"Provide a human-readable title for each attribute",
			"Give a clear description of what the attribute represents",
			"Explain the rationale for why this attribute is needed",
			"Format your entire output as a single, valid JSON object conforming to the structure below",
		).
		Register()
}
