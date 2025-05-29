package builtin

import (
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

// Register the processor with the registry
func init() {
	processor.NewBuilder("speech_act").
		WithStruct(&SpeechActResult{}).
		WithRole("You are an expert at identifying distinct speech acts within a text").
		WithObjective("Analyze the provided text and identify all distinct speech acts (like questions, requests, statements, greetings, etc.). For each identified speech act, provide its category, complexity, and relevant keywords").
		WithInstructions(
			"Read the text and identify each separate speech act - a single sentence might contain multiple speech acts",
			"For each speech act, determine its category (e.g., informational, question, request, command, greeting, confirmation)",
			"For each speech act, rate its complexity on a scale from 1.0 (very simple) to 10.0 (very complex)",
			"For each speech act, extract up to 3 relevant keywords",
			"Ensure the 'keywords' field for each speech act is a JSON array of strings",
			"If no relevant keywords are found for a specific speech act, use an empty array []",
		).
		Register()
}
