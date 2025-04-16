package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// KeywordResult contains the extracted keywords
type KeywordResult struct {
	Keywords      []string `json:"keywords"`
	Categories    []string `json:"categories,omitempty"`
	ProcessorType string   `json:"processor_type"`
}

// KeywordPrompt is a prompt generator for keyword extraction
type KeywordPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *KeywordPrompt) GeneratePrompt(_ context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	// Create an empty instance to generate the structure
	exampleResult := &KeywordResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert keyword extraction and categorization tool that ONLY outputs valid JSON.

**Input Text:**
%s

**Instructions:**
1. Extract the 5-10 most important keywords from the text.
2. Identify 2-3 categories that best describe the text.
3. Format your entire output as a single, valid JSON object.
4. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil // Use generated JSON example
}

func main() {
	// Get API Key from GEMINI_API_KEY environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable not set")
	}

	// Restore original LLM provider initialization using env var
	config := llm.Config{
		APIKey:      apiKey,
		Model:       "gemini-2.0-flash",
		MaxTokens:   1024,
		Temperature: 0.2,
	}

	provider, err := llm.NewProvider(llm.Google, config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Register the keyword processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"keyword",        // name
		[]string{"text"}, // contentTypes
		&KeywordResult{}, // resultStruct
		&KeywordPrompt{}, // promptGenerator
		nil,              // no custom initialization needed
		false,            // No struct validation needed by default
	)

	// Input text 1
	text1 := "Artificial intelligence (AI) is intelligence demonstrated by machines, " +
		"as opposed to intelligence displayed by humans or other animals. Example tasks in which " +
		"AI is applied include speech recognition, computer vision, translation between natural " +
		"languages, and other mappings of inputs. AI applications include advanced web search engines, " +
		"recommendation systems, understanding human speech, self-driving cars, automated decision-making, " +
		"and competing at the highest level in strategic game systems."

	// Input text 2
	text2 := "Machine learning is a subset of AI focused on training models to improve with experience."

	// Create ProcessItems
	item1 := data.NewTextProcessItem("example-1", text1, nil)
	item2 := data.NewTextProcessItem("example-2", text2, nil)

	// Get processor from registry
	// Use an empty Options struct for defaults
	keywordProcessor, err := processor.Create("keyword", provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor from registry: %v", err)
	}

	// Process the first item
	result1, err := keywordProcessor.Process(context.Background(), item1)
	if err != nil {
		log.Fatalf("Processing item 1 failed: %v", err)
	}

	// Print the first result
	fmt.Println("--- Result for Item 1 (ID: example-1) ---")
	if procInfo, ok := result1.ProcessingInfo["keyword"]; ok {
		jsonData, _ := json.MarshalIndent(procInfo, "", "  ")
		fmt.Println(string(jsonData))
	}

	// Process the second item
	result2, err := keywordProcessor.Process(context.Background(), item2)
	if err != nil {
		log.Fatalf("Processing item 2 failed: %v", err)
	}

	// Print the second result
	fmt.Println("\n--- Result for Item 2 (ID: example-2) ---")
	if procInfo, ok := result2.ProcessingInfo["keyword"]; ok {
		jsonData, _ := json.MarshalIndent(procInfo, "", "  ")
		fmt.Println(string(jsonData))
	}
}
