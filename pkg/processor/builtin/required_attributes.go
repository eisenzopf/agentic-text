package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
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

// RegisterDebugHandler wraps the processor with debugging
func RegisterDebugHandler() {
	// Register a custom initializer for the required_attributes processor
	// This adds debugging to help track what's happening
	processor.Register("required_attributes", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
		// Create a modified response handler that includes debugging
		responseHandler := &DebugResponseHandler{
			BaseHandler: &processor.BaseResponseHandler{
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

		// Create the processor using the proper API instead of accessing fields directly
		proc := processor.NewBaseProcessor(
			"required_attributes",
			[]string{"text", "json"},
			client,
			nil, // No pre-processor
			promptGen,
			responseHandler,
			options,
		)

		return proc, nil
	})
}

// DebugResponseHandler adds debugging to the processor
type DebugResponseHandler struct {
	BaseHandler *processor.BaseResponseHandler
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
								FieldName:   processor.GetStringValue(attrMap, "field_name"),
								Title:       processor.GetStringValue(attrMap, "title"),
								Description: processor.GetStringValue(attrMap, "description"),
								Rationale:   processor.GetStringValue(attrMap, "rationale"),
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
	handler := processor.NewResponseHandler("required_attributes", &RequiredAttributesResult{})

	// Process the response
	result, err := handler.AutoProcessResponse(ctx, text, responseData)

	// Debug the processed result
	fmt.Printf("DEBUG - Processed Result: %+v\n", result)

	return result, err
}

// DebugPromptGenerator adds debugging to the prompt generator
type DebugPromptGenerator struct {
	BaseGenerator processor.PromptGenerator
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
	processor.RegisterGenericProcessor(
		"required_attributes",       // name
		[]string{"text", "json"},    // contentTypes
		&RequiredAttributesResult{}, // resultStruct
		&RequiredAttributesPrompt{}, // promptGenerator
		nil,                         // no custom initialization needed
		map[string]interface{}{ // validation options
			"field_name": "attributes",
		},
	)

	// Register a debug version to help diagnose issues
	RegisterDebugHandler()
}
