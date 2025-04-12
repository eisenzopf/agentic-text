package easy

import (
	"encoding/json"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// ProcessText is a one-line function to process text with a specified processor type
func ProcessText(text, processorType string) (map[string]interface{}, error) {
	wrapper, err := New(processorType)
	if err != nil {
		return nil, err
	}

	return wrapper.Process(text)
}

// ProcessTextWithConfig processes text with a specified processor type and custom configuration
func ProcessTextWithConfig(text, processorType string, config *Config) (map[string]interface{}, error) {
	wrapper, err := NewWithConfig(processorType, config)
	if err != nil {
		return nil, err
	}

	return wrapper.Process(text)
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

	return wrapper.ProcessBatch(texts, concurrency)
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
