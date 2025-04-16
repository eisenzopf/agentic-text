package builtin

import (
	"context"
	"fmt"

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

// IntentPrompt is a prompt generator for intent analysis.
type IntentPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *IntentPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	// Provide an example with one intent item
	exampleResult := &IntentResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`You are a helpful AI assistant specializing in classifying customer service conversations. Your task is to analyze a provided conversation transcript and identify *all* distinct customer intents expressed.

**Input:** You will receive a conversation transcript as text. If the input appears to be in JSON format, focus on the text content and ignore the JSON structure.

**Output:** You will return a JSON object containing a list named "intents". Each item in the list should represent one distinct customer intent identified in the text, following this example structure:
%s

**Important Instructions and Constraints:**

1.  **Identify All Intents:** List every distinct reason the customer appears to be contacting support. If multiple intents are present, list them all.
2.  **Conciseness:** For each intent, keep the "label_name" to 2-3 words (Title Case) and the "description" brief and to the point (1-2 sentences).
3.  **JSON Format:** The output *must* be a valid JSON object containing the "intents" list. Do not include any extra text, explanations, or apologies outside of the JSON object. Only the JSON object should be returned.
4.  **Specificity:** Be as specific as possible in the description for each intent. Don't just say "billing issue." Say "The customer is disputing a charge on their latest bill."
5.  **Do not hallucinate information.** Base the classification solely on the provided transcript. Do not invent details.
6.  **Do not respond in a conversational manner.** Your entire response should be only the requested JSON.

Conversation Transcript:
%s`, jsonExample, text), nil
}

// Register the processor with the registry
func init() {
	// Register the intent processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"intent",                 // name
		[]string{"text", "json"}, // contentTypes
		&IntentResult{},          // resultStruct (defines the output structure)
		&IntentPrompt{},          // promptGenerator
		nil,                      // no custom initialization needed
		false,                    // No struct validation needed by default (defaults handled by items)
	)
}
