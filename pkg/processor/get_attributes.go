package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

	// Create client from provider
	client := llm.NewProviderClient(provider)

	// Create and embed base processor - support both text and json content types
	base := NewBaseProcessor("attributes", []string{"text", "json"}, client, nil, p, p, options)
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

// Helper function to safely get string values from interface maps
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// Helper function to safely get float values from interface maps
func getFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *AttributeProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	// Check if responseData is a string, which can happen with some providers
	if strResponse, ok := responseData.(string); ok {
		// Remove markdown code block if present
		cleanResponse := strResponse
		if strings.HasPrefix(cleanResponse, "```json") && strings.HasSuffix(cleanResponse, "```") {
			// Extract content between ```json and ```
			cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
			cleanResponse = strings.TrimSuffix(cleanResponse, "```")
			cleanResponse = strings.TrimSpace(cleanResponse)
		} else if strings.HasPrefix(cleanResponse, "```") && strings.HasSuffix(cleanResponse, "```") {
			// Extract content between ``` and ```
			cleanResponse = strings.TrimPrefix(cleanResponse, "```")
			cleanResponse = strings.TrimSuffix(cleanResponse, "```")
			cleanResponse = strings.TrimSpace(cleanResponse)
		}

		// Try to parse the string as JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(cleanResponse), &data); err != nil {
			// If parsing fails, wrap it as a response with empty attributes
			return map[string]interface{}{
				"attributes":     []Attribute{},
				"response":       strResponse,
				"processor_type": "get_attributes",
			}, nil
		}
		// If parsing succeeds, use the parsed data
		responseData = data
	}

	// Convert the response data to map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format: %T", responseData)
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
					FieldName:   getStringValue(attrMap, "field_name"),
					Value:       getStringValue(attrMap, "value"),
					Confidence:  getFloatValue(attrMap, "confidence"),
					Explanation: getStringValue(attrMap, "explanation"),
				}
				attributes = append(attributes, attribute)
			}
		}
	}

	// Create result map
	resultMap := map[string]interface{}{
		"attributes":     attributes,
		"processor_type": "get_attributes",
	}

	// Add debug info back if it existed
	if debugInfo != nil {
		resultMap["debug"] = debugInfo
	}

	return resultMap, nil
}

// Register the processor with the registry
func init() {
	Register("get_attributes", func(provider llm.Provider, options Options) (Processor, error) {
		return NewAttributeProcessor(provider, options)
	})
}
