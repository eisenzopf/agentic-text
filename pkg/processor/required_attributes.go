package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
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
		var attributesRaw []interface{}

		// Handle different ways the LLM might return data
		switch v := val.(type) {
		case []interface{}:
			// Direct array of attributes
			attributesRaw = v
		case map[string]interface{}:
			// Attributes in a nested "attributes" field
			if attrs, ok := v["attributes"].([]interface{}); ok {
				attributesRaw = attrs
			} else {
				return defaultAttr
			}
		default:
			return defaultAttr
		}

		// If no attributes, return default
		if len(attributesRaw) == 0 {
			return defaultAttr
		}

		// Process each attribute to ensure it has the right structure
		validAttributes := make([]AttributeDefinition, 0, len(attributesRaw))
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

			// Create a valid AttributeDefinition
			attr := AttributeDefinition{
				FieldName:   fieldName,
				Title:       GetStringValue(attrMap, "title"),
				Description: GetStringValue(attrMap, "description"),
				Rationale:   GetStringValue(attrMap, "rationale"),
			}

			// Add the validated attribute
			validAttributes = append(validAttributes, attr)
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
	// Generate example JSON from the result struct
	exampleResult := &RequiredAttributesResult{}
	jsonExample := GenerateJSONExample(exampleResult)

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

// RegisterDebugHandler wraps the processor with debugging
func RegisterDebugHandler() {
	// Register a custom initializer for the required_attributes processor
	// This adds debugging to help track what's happening
	Register("required_attributes", func(provider llm.Provider, options Options) (Processor, error) {
		// Create a modified response handler that includes debugging
		responseHandler := &DebugResponseHandler{
			BaseHandler: &BaseResponseHandler{
				ProcessorType: "required_attributes",
				ResultStruct:  &RequiredAttributesResult{},
			},
		}

		// Create the prompt generator
		promptGen := &DebugPromptGenerator{
			BaseGenerator: &RequiredAttributesPrompt{},
		}

		// Create client from provider
		client := llm.NewProviderClient(provider)

		// Create the processor with our debug handler
		proc := &BaseProcessor{
			name:            "required_attributes",
			contentTypes:    []string{"text", "json"},
			llmClient:       client,
			promptGenerator: promptGen,
			responseHandler: responseHandler,
			options:         options,
		}

		return proc, nil
	})
}

// DebugResponseHandler adds debugging to the processor
type DebugResponseHandler struct {
	BaseHandler *BaseResponseHandler
}

// HandleResponse implements ResponseHandler interface
func (h *DebugResponseHandler) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	// Debug the raw response
	fmt.Printf("DEBUG - Raw LLM Response: %+v\n", responseData)

	// Handle string responses with markdown code blocks
	if strResponse, ok := responseData.(string); ok {
		// Check if it's wrapped in markdown code blocks
		if strings.HasPrefix(strResponse, "```") {
			fmt.Println("DEBUG - Detected markdown code block, cleaning...")
			// Remove markdown formatting
			strResponse = strings.TrimPrefix(strResponse, "```json")
			strResponse = strings.TrimPrefix(strResponse, "```")
			endIndex := strings.LastIndex(strResponse, "```")
			if endIndex != -1 {
				strResponse = strResponse[:endIndex]
			}
			strResponse = strings.TrimSpace(strResponse)

			// Parse the cleaned JSON
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(strResponse), &jsonData); err == nil {
				fmt.Printf("DEBUG - Cleaned JSON: %+v\n", jsonData)

				// Check if we have attributes in the JSON
				if attrsRaw, ok := jsonData["attributes"].([]interface{}); ok && len(attrsRaw) > 0 {
					// Create a new result with the attributes
					result := &RequiredAttributesResult{
						ProcessorType: "required_attributes",
						Attributes:    make([]AttributeDefinition, 0, len(attrsRaw)),
					}

					// Convert each attribute
					for _, attrRaw := range attrsRaw {
						if attrMap, ok := attrRaw.(map[string]interface{}); ok {
							// Create a new attribute
							attr := AttributeDefinition{
								FieldName:   GetStringValue(attrMap, "field_name"),
								Title:       GetStringValue(attrMap, "title"),
								Description: GetStringValue(attrMap, "description"),
								Rationale:   GetStringValue(attrMap, "rationale"),
							}

							// Add it to the result if it has a field name
							if attr.FieldName != "" {
								result.Attributes = append(result.Attributes, attr)
							}
						}
					}

					// Return the result if we have attributes
					if len(result.Attributes) > 0 {
						fmt.Printf("DEBUG - Created result: %+v\n", result)
						return result, nil
					}
				}

				// If we got here, use the cleaned JSON as the response
				responseData = jsonData
			} else {
				fmt.Printf("DEBUG - Error parsing cleaned JSON: %v\n", err)
			}
		}
	}

	// Create a standard handler for fallback
	handler := NewResponseHandler("required_attributes", &RequiredAttributesResult{})

	// Process the response
	result, err := handler.AutoProcessResponse(ctx, text, responseData)

	// Debug the processed result
	fmt.Printf("DEBUG - Processed Result: %+v\n", result)

	return result, err
}

// DebugPromptGenerator adds debugging to the prompt generator
type DebugPromptGenerator struct {
	BaseGenerator PromptGenerator
}

// GeneratePrompt implements PromptGenerator interface
func (p *DebugPromptGenerator) GeneratePrompt(ctx context.Context, text string) (string, error) {
	prompt, err := p.BaseGenerator.GeneratePrompt(ctx, text)
	if err != nil {
		return "", err
	}

	// Print the prompt for debugging
	fmt.Println("DEBUG - LLM Prompt:")
	fmt.Println("====================================")
	fmt.Println(prompt)
	fmt.Println("====================================")

	return prompt, nil
}

// init registers the standard and debug processors
func init() {
	// Register the standard processor using the new validation approach
	RegisterGenericProcessor(
		"required_attributes",       // name
		[]string{"text", "json"},    // contentTypes
		&RequiredAttributesResult{}, // resultStruct
		&RequiredAttributesPrompt{}, // promptGenerator
		nil,                         // no custom initialization needed
		map[string]interface{}{ // validation options
			"field_name": "attributes",
			"default_value": []AttributeDefinition{
				{
					FieldName:   "unknown",
					Title:       "Unknown",
					Description: "Unable to determine required attributes from the response",
					Rationale:   "The response did not contain valid attribute definitions",
				},
			},
		},
	)

	// Register a debug version to help diagnose issues
	RegisterDebugHandler()
}
