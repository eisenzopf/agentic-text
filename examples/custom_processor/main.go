package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// KeywordResult contains the extracted keywords
type KeywordResult struct {
	Keywords      []string       `json:"keywords"`
	Categories    []string       `json:"categories,omitempty"`
	Frequencies   map[string]int `json:"frequencies,omitempty"`
	ProcessorType string         `json:"processor_type"`
}

// KeywordPrompt is a prompt generator for keyword extraction
type KeywordPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *KeywordPrompt) GeneratePrompt(_ context.Context, text string) (string, error) {
	return fmt.Sprintf(`**Role:** You are an expert keyword extraction tool that ONLY outputs valid JSON.

**Input Text:**
%s

**Instructions:**
1. Extract the 5-10 most important keywords from the text
2. Identify 2-3 categories that best describe the text
3. Count the frequency of each keyword in the text
4. Format your entire output as a single, valid JSON object
5. *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text. ***

**Required JSON Output Structure:**
{
  "keywords": ["...", "..."],     // Array of 5-10 important keywords
  "categories": ["...", "..."],   // Array of 2-3 categories
  "frequencies": {                // Map of keyword frequencies
    "keyword1": 3,
    "keyword2": 2
  }
}`, text), nil
}

func main() {
	// Initialize LLM provider
	config := llm.Config{
		APIKey:      "your-api-key",
		Model:       "gemini-pro",
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
	)

	// Input text
	text := "Artificial intelligence (AI) is intelligence demonstrated by machines, " +
		"as opposed to intelligence displayed by humans or other animals. Example tasks in which " +
		"AI is applied include speech recognition, computer vision, translation between natural " +
		"languages, and other mappings of inputs. AI applications include advanced web search engines, " +
		"recommendation systems, understanding human speech, self-driving cars, automated decision-making, " +
		"and competing at the highest level in strategic game systems."

	// Create a ProcessItem from the text
	item := data.NewTextProcessItem("example-1", text, nil)

	// Get processor from registry
	keywordProcessor, err := processor.Create("keyword", provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor from registry: %v", err)
	}

	// Process the item
	result, err := keywordProcessor.Process(context.Background(), item)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print the result
	fmt.Println("Keyword Processor Result:")
	if procInfo, ok := result.ProcessingInfo["keyword"]; ok {
		jsonData, _ := json.MarshalIndent(procInfo, "", "  ")
		fmt.Println(string(jsonData))
	}

	// Create a new ProcessItem for the second example
	secondItem := data.NewTextProcessItem("example-2",
		"Machine learning is a subset of AI focused on training models to improve with experience.",
		nil)

	// Process the second item
	regResult, err := keywordProcessor.Process(context.Background(), secondItem)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print JSON result
	fmt.Println("\nSecond Example Result:")
	if procInfo, ok := regResult.ProcessingInfo["keyword"]; ok {
		jsonData, _ := json.MarshalIndent(procInfo, "", "  ")
		fmt.Println(string(jsonData))
	}
}
