package builtin

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// KeywordResult contains the keyword extraction results
type KeywordResult struct {
	// Keywords is an array of extracted keywords
	Keywords []Keyword `json:"keywords,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// DefaultValues returns the default values for this result type
func (r *KeywordResult) DefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"keywords": []Keyword{},
	}
}

// Keyword represents a single extracted keyword
type Keyword struct {
	// Term is the extracted keyword term
	Term string `json:"term"`
	// Relevance is the relevance score from 0.0 to 1.0
	Relevance float64 `json:"relevance"`
	// Category is the category of the keyword (e.g., "topic", "person", "location")
	Category string `json:"category"`
}

// KeywordPrompt is a prompt generator for keyword extraction
type KeywordPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *KeywordPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	exampleResult := &KeywordResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert at extracting important keywords from text.

**Objective:** Analyze the provided text and extract the most meaningful keywords.

**Input Text:**
%s

**Instructions:**
1. Carefully read and interpret the Input Text.
2. Extract the most important keywords or key phrases that represent the main topics, following the structure below.
3. For each keyword, provide:
   - The keyword term
   - A relevance score (0.0 to 1.0) indicating how central the keyword is to the content
   - A category for the keyword (e.g., "topic", "person", "location", "concept", "organization")
4. Format your entire output as a single, valid JSON object.
5. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil
}

// Register the processor with the registry
func init() {
	// Register the keyword processor using generic processor registration with validation
	processor.RegisterGenericProcessor(
		"keyword_extraction", // name
		[]string{"text"},     // contentTypes
		&KeywordResult{},     // resultStruct
		&KeywordPrompt{},     // promptGenerator
		nil,                  // no custom initialization needed
		map[string]interface{}{ // validation options
			"field_name":    "keywords",
			"default_value": []Keyword{},
		},
	)
}
