package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// IntentProcessor analyzes the intent of customer service conversations
type IntentProcessor struct {
	// Embed BaseProcessor to inherit all methods
	BaseProcessor
}

// IntentResult contains the intent classification results
type IntentResult struct {
	// LabelName is a natural language label describing the customer's primary intent (title case)
	LabelName string `json:"label_name"`
	// Label is a machine-readable version of LabelName (snake_case)
	Label string `json:"label"`
	// Description is a concise description of the customer's primary intent
	Description string `json:"description"`
}

// NewIntentProcessor creates a new intent processor
func NewIntentProcessor(provider llm.Provider, options Options) (*IntentProcessor, error) {
	p := &IntentProcessor{}

	// Create client from provider
	client := llm.NewProviderClient(provider)

	// Create and embed base processor - support both text and json content types
	base := NewBaseProcessor("intent", []string{"text", "json"}, client, nil, p, p, options)
	p.BaseProcessor = *base

	return p, nil
}

// GeneratePrompt implements PromptGenerator interface - generates the intent analysis prompt
func (p *IntentProcessor) GeneratePrompt(ctx context.Context, text string) (string, error) {
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

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *IntentProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
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
			// If parsing fails, wrap it as a response
			return map[string]interface{}{
				"label_name":     "Unclear Intent",
				"label":          "unclear_intent",
				"description":    "The conversation transcript is unclear or does not contain a discernible customer service request.",
				"response":       strResponse,
				"processor_type": "intent",
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

	// Check if we got a non-JSON response wrapped in a "response" field
	if response, exists := data["response"]; exists && len(data) <= 2 { // data has only response and maybe debug
		// This is a fallback case where the LLM didn't produce valid JSON
		resultMap := map[string]interface{}{
			"label_name":     "Unclear Intent",
			"label":          "unclear_intent",
			"description":    "The conversation transcript is unclear or does not contain a discernible customer service request.",
			"response":       response,
			"processor_type": "intent",
		}

		// Add debug info back if it existed
		if debugInfo != nil {
			resultMap["debug"] = debugInfo
		}

		return resultMap, nil
	}

	// Extract intent fields with defaults
	labelName, _ := data["label_name"].(string)
	label, _ := data["label"].(string)
	description, _ := data["description"].(string)

	// Set defaults if missing
	if labelName == "" || label == "" {
		labelName = "Unclear Intent"
		label = "unclear_intent"
		if description == "" {
			description = "The conversation transcript is unclear or does not contain a discernible customer service request."
		}
	}

	// Create result map with intent data
	resultMap := map[string]interface{}{
		"label_name":     labelName,
		"label":          label,
		"description":    description,
		"processor_type": "intent",
	}

	// Add debug info back if it existed
	if debugInfo != nil {
		resultMap["debug"] = debugInfo
	}

	return resultMap, nil
}

// Register the processor with the registry
func init() {
	Register("intent", func(provider llm.Provider, options Options) (Processor, error) {
		return NewIntentProcessor(provider, options)
	})
}
