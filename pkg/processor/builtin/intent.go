package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// IntentItem represents a single identified customer intent.
type IntentItem struct {
	// LabelName is a natural language label describing the customer's intent (title case, 2-3 words)
	LabelName string `json:"label_name" default:"Unclear Intent"`
	// Label is a machine-readable version of LabelName (snake_case)
	Label string `json:"label" default:"unclear_intent"`
	// Description is a concise description of the customer's intent (1-2 sentences)
	Description string `json:"description" default:"The conversation transcript is unclear or does not contain a discernible customer service request."`
}

// IntentResult contains a list of identified customer intents.
type IntentResult struct {
	// Intents is a list of intents identified in the conversation.
	Intents []IntentItem `json:"intents"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("intent").
		WithStruct(&IntentResult{}).
		WithContentTypes("text", "json").
		WithRole("You are a helpful AI assistant specializing in classifying customer service conversations").
		WithObjective("Analyze a provided conversation transcript and identify *all* distinct customer intents expressed").
		WithInstructions(
			"Identify All Intents: List every distinct reason the customer appears to be contacting support",
			"If multiple intents are present, list them all",
			"Keep the 'label_name' to 2-3 words (Title Case) and the 'description' brief and to the point (1-2 sentences)",
			"Be as specific as possible in the description for each intent",
			"Don't just say 'billing issue.' Say 'The customer is disputing a charge on their latest bill.'",
			"Do not hallucinate information. Base the classification solely on the provided transcript",
		).
		WithCustomSection("Important Constraints", `
- Do not respond in a conversational manner
- Your entire response should be only the requested JSON
- If the input appears to be in JSON format, focus on the text content and ignore the JSON structure`).
		Register()
}
