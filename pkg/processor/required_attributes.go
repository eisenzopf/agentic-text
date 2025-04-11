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

// RequiredAttributesConfig holds configuration for the RequiredAttributesProcessor
type RequiredAttributesConfig struct {
	Questions          []string
	ExistingAttributes []string
}

// RequiredAttributesPromptGenerator generates prompts for attribute definition
type RequiredAttributesPromptGenerator struct {
	config RequiredAttributesConfig
}

// GeneratePrompt implements the PromptGenerator interface
func (pg *RequiredAttributesPromptGenerator) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Format questions for the prompt
	questionsText := ""
	for i, q := range pg.config.Questions {
		questionsText += fmt.Sprintf("%d. %s\n", i+1, q)
	}

	// Format existing attributes if provided
	existingText := ""
	if len(pg.config.ExistingAttributes) > 0 {
		existingText = "\nExisting attributes:\n"
		for _, attr := range pg.config.ExistingAttributes {
			existingText += fmt.Sprintf("- %s\n", attr)
		}
	}

	prompt := fmt.Sprintf(`We need to determine what data attributes are required to answer these questions:
%s
%s

Return a JSON object with this structure:
{
  "attributes": [
    {
      "field_name": str,  // Database field name in snake_case
      "title": str,       // Human readable title
      "description": str, // Detailed description of the attribute
      "rationale": str    // Why this attribute is needed for the questions
    }
  ]
}`, questionsText, existingText)

	return prompt, nil
}

// RequiredAttributesResponseHandler handles responses from the LLM
type RequiredAttributesResponseHandler struct{}

// HandleResponse implements the ResponseHandler interface
func (rh *RequiredAttributesResponseHandler) HandleResponse(ctx context.Context, text string, responseData interface{}) (*Result, error) {
	// Validate response format
	resultMap, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result format")
	}

	attributesRaw, ok := resultMap["attributes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("attributes field is not an array")
	}

	// Convert to strongly typed attributes
	attributes := make([]AttributeDefinition, 0, len(attributesRaw))
	for _, attrRaw := range attributesRaw {
		attrMap, ok := attrRaw.(map[string]interface{})
		if !ok {
			continue // Skip invalid entries
		}

		attr := AttributeDefinition{
			FieldName:   getString(attrMap, "field_name"),
			Title:       getString(attrMap, "title"),
			Description: getString(attrMap, "description"),
			Rationale:   getString(attrMap, "rationale"),
		}

		// Only add if field_name is valid
		if attr.FieldName != "" {
			attributes = append(attributes, attr)
		}
	}

	return &Result{
		Processed: text,
		Data:      attributes,
	}, nil
}

// RequiredAttributesProcessor generates required attributes based on questions
type RequiredAttributesProcessor struct {
	*BaseProcessor
	config RequiredAttributesConfig
}

// NewRequiredAttributesProcessor creates a new RequiredAttributesProcessor
func NewRequiredAttributesProcessor(provider llm.Provider, config RequiredAttributesConfig, options Options) *RequiredAttributesProcessor {
	promptGen := &RequiredAttributesPromptGenerator{config: config}
	respHandler := &RequiredAttributesResponseHandler{}

	baseProcessor := NewBaseProcessor(
		"required_attributes",
		provider,
		options,
		&DefaultPreProcessor{},
		promptGen,
		respHandler,
	)

	return &RequiredAttributesProcessor{
		BaseProcessor: baseProcessor,
		config:        config,
	}
}

// Process overrides the default process method to handle the empty text case
func (p *RequiredAttributesProcessor) Process(ctx context.Context, text string) (*Result, error) {
	// Validate request
	if len(p.config.Questions) == 0 {
		return nil, fmt.Errorf("questions are required")
	}

	// Use the base processor implementation
	return p.BaseProcessor.Process(ctx, text)
}

// GetAttributes returns the attribute definitions from the result
func GetAttributes(result *Result) ([]AttributeDefinition, error) {
	if result == nil {
		return nil, fmt.Errorf("result is nil")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	attributes, ok := result.Data.([]AttributeDefinition)
	if !ok {
		return nil, fmt.Errorf("result data is not of type []AttributeDefinition")
	}

	return attributes, nil
}

// Helper function to extract string values from maps
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
