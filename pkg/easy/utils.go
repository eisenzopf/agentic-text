package easy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// CleanLLMResponse ensures that we extract proper JSON from LLM responses
// This handles the case where we get responses with raw JSON inside
func CleanLLMResponse(response map[string]interface{}) map[string]interface{} {
	// If we have a "response" field that contains a JSON string, try to parse and use it
	if responseStr, ok := response["response"].(string); ok && strings.Contains(responseStr, "{") {
		// Extract JSON from the response if it's in a code block
		cleanResponse := responseStr

		// Try to find JSON between code blocks
		if strings.Contains(cleanResponse, "```") {
			// Handle ```json blocks
			if strings.Contains(cleanResponse, "```json") {
				parts := strings.Split(cleanResponse, "```json")
				if len(parts) > 1 {
					codeContent := parts[1]
					endPos := strings.Index(codeContent, "```")
					if endPos != -1 {
						cleanResponse = strings.TrimSpace(codeContent[:endPos])
					}
				}
			} else {
				// Try to extract content from generic code blocks
				parts := strings.Split(cleanResponse, "```")
				if len(parts) >= 3 { // At least one complete code block
					// Take the content of the first code block
					cleanResponse = strings.TrimSpace(parts[1])
				}
			}
		}

		// Try to parse as JSON
		var extractedData map[string]interface{}
		if err := json.Unmarshal([]byte(cleanResponse), &extractedData); err == nil {
			// Successfully parsed JSON from response field
			// Add the processor_type from the original result
			if procType, exists := response["processor_type"]; exists {
				extractedData["processor_type"] = procType
			}

			// Keep the original response for debugging if needed
			extractedData["original_response"] = responseStr

			return extractedData
		}
	}

	// Return the original response if no JSON could be extracted
	return response
}

// ProcessText is a one-line function to process text with a specified processor type
func ProcessText(text, processorType string) (map[string]interface{}, error) {
	wrapper, err := New(processorType)
	if err != nil {
		return nil, err
	}

	result, err := wrapper.Process(text)
	if err != nil {
		return nil, err
	}

	// Clean the response to handle JSON inside response field
	return CleanLLMResponse(result), nil
}

// ProcessTextWithConfig processes text with a specified processor type and custom configuration
func ProcessTextWithConfig(text, processorType string, config *Config) (map[string]interface{}, error) {
	wrapper, err := NewWithConfig(processorType, config)
	if err != nil {
		return nil, err
	}

	result, err := wrapper.Process(text)
	if err != nil {
		return nil, err
	}

	// Clean the response to handle JSON inside response field
	return CleanLLMResponse(result), nil
}

// ProcessBatchText processes a batch of texts with a specified processor type
func ProcessBatchText(texts []string, processorType string, concurrency int) ([]map[string]interface{}, error) {
	if concurrency <= 0 {
		concurrency = 2
	}

	wrapper, err := New(processorType)
	if err != nil {
		return nil, err
	}

	results, err := wrapper.ProcessBatch(texts, concurrency)
	if err != nil {
		return nil, err
	}

	// Clean the responses to handle JSON inside response fields
	cleanedResults := make([]map[string]interface{}, len(results))
	for i, result := range results {
		cleanedResults[i] = CleanLLMResponse(result)
	}

	return cleanedResults, nil
}

// ProcessJSON takes a JSON string, processes it with the specified processor, and returns a result
func ProcessJSON(jsonStr, processorType string) (map[string]interface{}, error) {
	// Parse JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Extract text field if exists
	var text string
	if textVal, ok := jsonData["text"].(string); ok {
		text = textVal
	} else if contentVal, ok := jsonData["content"].(string); ok {
		text = contentVal
	} else {
		return nil, fmt.Errorf("JSON must contain a 'text' or 'content' field as string")
	}

	// Process text
	wrapper, err := New(processorType)
	if err != nil {
		return nil, err
	}

	return wrapper.Process(text)
}

// ListAvailableProcessors returns a list of all registered processor types
func ListAvailableProcessors() []string {
	// Import the processor package to ensure all processors are registered
	return processor.ListProcessors()
}

// PrettyPrint formats a result map as a readable JSON string
func PrettyPrint(result map[string]interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Sentiment analyzes the sentiment of the given text
func Sentiment(text string) (map[string]interface{}, error) {
	return ProcessText(text, "sentiment")
}

// Intent analyzes the intent in the given text
func Intent(text string) (map[string]interface{}, error) {
	return ProcessText(text, "intent")
}

// RequiredAttributes identifies required attributes in the given text
func RequiredAttributes(text string) (map[string]interface{}, error) {
	return ProcessText(text, "required_attributes")
}

// GetAttributes identifies and extracts attributes from the given text
func GetAttributes(text string) (map[string]interface{}, error) {
	return ProcessText(text, "get_attributes")
}
