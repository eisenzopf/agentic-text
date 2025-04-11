package processor

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// AttributeProcessor extracts attributes from text
type AttributeProcessor struct {
	// Embed BaseProcessor to inherit all methods
	BaseProcessor
}

// AttributeResult contains the extracted attributes
type AttributeResult struct {
	// Attributes is an array of extracted attributes
	Attributes []Attribute `json:"attributes"`
	// Confidence is the confidence level (0.0 to 1.0)
	Confidence float64 `json:"confidence"`
	// Explanation provides context for the extracted attributes
	Explanation string `json:"explanation"`
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

// NewAttributeProcessor creates a new attribute processor
func NewAttributeProcessor(provider llm.Provider, options Options) (*AttributeProcessor, error) {
	p := &AttributeProcessor{}

	// Create and embed base processor
	base := NewBaseProcessor("attributes", provider, options, nil, p, p)
	p.BaseProcessor = *base

	return p, nil
}

// GeneratePrompt implements PromptGenerator interface - generates the attribute extraction prompt
func (p *AttributeProcessor) GeneratePrompt(ctx context.Context, text string) (string, error) {
	return fmt.Sprintf(`**Role:** You are an expert at extracting structured information from text.

**Objective:** Analyze the provided text and extract relevant attributes and their values.

**Input Text:**
%s

**Instructions:**
1. Carefully read and interpret the Input Text.
2. Extract any relevant attributes and their values.
3. For each attribute, provide:
   - A field name (in snake_case)
   - The extracted value
   - A confidence score (0.0 to 1.0)
   - A brief explanation
4. Assign an overall confidence score for the extraction.
5. Provide a brief overall explanation of how the attributes were determined.
6. Format your entire output as a single, valid JSON object.
7. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

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
  ],
  "confidence": 0.0,
  "explanation": "..."
}`, text), nil
}

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *AttributeProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (*Result, error) {
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
	var attributes []Attribute
	if attrs, ok := data["attributes"].([]interface{}); ok {
		for _, attr := range attrs {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				attribute := Attribute{
					FieldName:   getString(attrMap, "field_name"),
					Value:       getString(attrMap, "value"),
					Confidence:  getFloat(attrMap, "confidence"),
					Explanation: getString(attrMap, "explanation"),
				}
				attributes = append(attributes, attribute)
			}
		}
	}

	// Extract overall confidence and explanation
	confidence, _ := data["confidence"].(float64)
	explanation, _ := data["explanation"].(string)

	// Create result map
	resultMap := map[string]interface{}{
		"attributes":  attributes,
		"confidence":  confidence,
		"explanation": explanation,
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

// Helper function to safely get float values from interface maps
func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

// Register the processor with the registry
func init() {
	Register("get_attributes", func(provider llm.Provider, options Options) (Processor, error) {
		return NewAttributeProcessor(provider, options)
	})
}
