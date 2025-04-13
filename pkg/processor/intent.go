package processor

import (
	"context"
	"fmt"
)

// IntentResult contains the intent classification results
type IntentResult struct {
	// LabelName is a natural language label describing the customer's primary intent (title case)
	LabelName string `json:"label_name"`
	// Label is a machine-readable version of LabelName (snake_case)
	Label string `json:"label"`
	// Description is a concise description of the customer's primary intent
	Description string `json:"description"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// IntentPrompt is a prompt generator for intent analysis
type IntentPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *IntentPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	return fmt.Sprintf(`You are a helpful AI assistant specializing in classifying customer service conversations. Your task is to analyze a provided conversation transcript and determine the customer's *primary* intent for contacting customer service. Focus on the *main reason* the customer initiated the interaction, even if other topics are briefly mentioned.

**Input:** You will receive a conversation transcript as text. If the input appears to be in JSON format, focus on the text content and ignore the JSON structure.

**Output:** You will return a JSON object with the following *exact* keys and data types:

* **"label_name"**: (string) A natural language label describing the customer's primary intent. This label should be 2-3 words *maximum*. Use title case (e.g., "Update Address", "Cancel Order").
* **"label"**: (string) A lowercase version of "label_name", with underscores replacing spaces (e.g., "update_address", "cancel_order"). This should be machine-readable.
* **"description"**: (string) A concise description (1-2 sentences) of the customer's primary intent. Explain the *specific* problem or request the customer is making.

**Important Instructions and Constraints:**

1. **Primary Intent Focus:** Identify the *single, most important* reason the customer contacted support. Ignore minor side issues if they are not the core reason for the interaction.
2. **Conciseness:** Keep the "label_name" to 2-3 words and the "description" brief and to the point.
3. **JSON Format:** The output *must* be valid JSON. Do not include any extra text, explanations, or apologies outside of the JSON object. Only the JSON object should be returned.
4. **Specificity:** Be as specific as possible in the description. Don't just say "billing issue." Say "The customer is disputing a charge on their latest bill."
5. **Do not hallucinate information.** Base the classification solely on the provided transcript. Do not invent details.
6. **Do not respond in a conversational manner.** Your entire response should be only the requested json.

Conversation Transcript:
%s`, text), nil
}

// Register the processor with the registry
func init() {
	// Register the intent processor using the generic processor registration
	RegisterGenericProcessor(
		"intent",                 // name
		[]string{"text", "json"}, // contentTypes
		&IntentResult{},          // resultStruct
		&IntentPrompt{},          // promptGenerator
		nil,                      // no custom initialization needed
	)
}
