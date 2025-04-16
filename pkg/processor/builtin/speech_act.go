package builtin

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// SpeechActItem represents a single identified speech act within the text.
type SpeechActItem struct {
	// Category of the speech act (informational, question, request, greeting, command, etc.)
	Category string `json:"category" default:"request"`
	// Complexity score from 1-10 related to this specific speech act
	Complexity float64 `json:"complexity" default:"1.0"`
	// Keywords relevant to this specific speech act
	Keywords []string `json:"keywords"`
}

// SpeechActResult contains a list of identified speech acts.
type SpeechActResult struct {
	// SpeechActs is a list of speech acts identified in the text.
	SpeechActs []SpeechActItem `json:"speech_acts"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// SpeechActPrompt is a prompt generator for the speech act processor.
type SpeechActPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *SpeechActPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	// We need an example with at least one item to show the structure
	exampleResult := &SpeechActResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert at identifying distinct speech acts within a text.

**Objective:** Analyze the provided text and identify all distinct speech acts (like questions, requests, statements, greetings, etc.). For each identified speech act, provide its category, complexity, and relevant keywords.

**Input Text:**
%s

**Instructions:**
1. Read the text and identify each separate speech act. A single sentence might contain multiple speech acts.
2. For each speech act, determine its category (e.g., informational, question, request, command, greeting, confirmation).
3. For each speech act, rate its complexity on a scale from 1.0 (very simple) to 10.0 (very complex).
4. For each speech act, extract up to 3 relevant keywords.
5. Format your response as a valid JSON object containing a list named "speech_acts". Each item in the list should represent one identified speech act.
6. Ensure the 'keywords' field for each speech act is a JSON array of strings (e.g., ["word1", "word2"]). If no relevant keywords are found for a specific speech act, use an empty array [].
7. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure (Example):**
%s`, text, jsonExample), nil
}

// Register the processor with the registry
func init() {
	// Register the speech act processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"speech_act",       // name
		[]string{"text"},   // contentTypes
		&SpeechActResult{}, // resultStruct (defines the output structure)
		&SpeechActPrompt{}, // promptGenerator
		nil,                // no custom initialization needed
		false,              // No struct validation needed by default (defaults handled by items)
	)
}
