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

// KeywordProcessor extracts keywords from text
type KeywordProcessor struct {
	processor.BaseProcessor
}

// KeywordResult contains the extracted keywords
type KeywordResult struct {
	Keywords    []string       `json:"keywords"`
	Categories  []string       `json:"categories"`
	Frequencies map[string]int `json:"frequencies"`
}

// NewKeywordProcessor creates a new keyword processor
func NewKeywordProcessor(provider llm.Provider, options processor.Options) (*KeywordProcessor, error) {
	p := &KeywordProcessor{}

	// Create client from provider
	client := llm.NewProviderClient(provider)

	// Pass the processor itself as the implementations for required interfaces
	base := processor.NewBaseProcessor("keyword", []string{"text"}, client, nil, p, p, options)
	p.BaseProcessor = *base

	return p, nil
}

// GeneratePrompt implements PromptGenerator interface
func (p *KeywordProcessor) GeneratePrompt(_ context.Context, text string) (string, error) {
	return fmt.Sprintf(`Extract the most important keywords from the following text:
Text: %s

Respond with a JSON object containing:
- "keywords": An array of the 5-10 most important keywords
- "categories": An array of 2-3 categories that best describe the text
- "frequencies": A map of how many times each keyword appears

Format your response as valid JSON.`, text), nil
}

// HandleResponse implements ResponseHandler interface
func (p *KeywordProcessor) HandleResponse(_ context.Context, text string, responseData interface{}) (interface{}, error) {
	// Convert the response data to KeywordResult
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format")
	}

	// Extract keywords
	var keywords []string
	if keywordsData, ok := data["keywords"].([]interface{}); ok {
		for _, k := range keywordsData {
			if keyword, ok := k.(string); ok {
				keywords = append(keywords, keyword)
			}
		}
	}

	// Extract categories
	var categories []string
	if categoriesData, ok := data["categories"].([]interface{}); ok {
		for _, c := range categoriesData {
			if category, ok := c.(string); ok {
				categories = append(categories, category)
			}
		}
	}

	// Extract frequencies
	frequencies := make(map[string]int)
	if freqData, ok := data["frequencies"].(map[string]interface{}); ok {
		for k, v := range freqData {
			if freq, ok := v.(float64); ok {
				frequencies[k] = int(freq)
			}
		}
	}

	// Create keyword result
	keywordResult := map[string]interface{}{
		"keywords":    keywords,
		"categories":  categories,
		"frequencies": frequencies,
	}

	return keywordResult, nil
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

	// Create our custom processor
	keywordProcessor, err := NewKeywordProcessor(provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to create processor: %v", err)
	}

	// Register our processor with the registry
	processor.Register("keyword", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
		return NewKeywordProcessor(provider, options)
	})

	// Input text
	text := "Artificial intelligence (AI) is intelligence demonstrated by machines, " +
		"as opposed to intelligence displayed by humans or other animals. Example tasks in which " +
		"AI is applied include speech recognition, computer vision, translation between natural " +
		"languages, and other mappings of inputs. AI applications include advanced web search engines, " +
		"recommendation systems, understanding human speech, self-driving cars, automated decision-making, " +
		"and competing at the highest level in strategic game systems."

	// Create a ProcessItem from the text
	item := data.NewTextProcessItem("example-1", text, nil)

	// Process the item
	result, err := keywordProcessor.Process(context.Background(), item)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print the result
	fmt.Println("Custom Processor Result:")

	// Get the keyword data from the ProcessingInfo
	if procInfo, ok := result.ProcessingInfo["keyword"]; ok {
		if keywordData, ok := procInfo.(map[string]interface{}); ok {
			fmt.Printf("Keywords: %v\n", keywordData["keywords"])
			fmt.Printf("Categories: %v\n", keywordData["categories"])
			fmt.Println("Frequencies:")
			if freqs, ok := keywordData["frequencies"].(map[string]interface{}); ok {
				for k, v := range freqs {
					fmt.Printf("  %s: %v\n", k, v)
				}
			}
		}
	}

	// Verify that our processor was registered correctly
	regProcessor, err := processor.Create("keyword", provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor from registry: %v", err)
	}

	// Create a new ProcessItem for the second example
	secondItem := data.NewTextProcessItem("example-2",
		"Machine learning is a subset of AI focused on training models to improve with experience.",
		nil)

	// Use the registered processor
	regResult, err := regProcessor.Process(context.Background(), secondItem)
	if err != nil {
		log.Fatalf("Processing with registered processor failed: %v", err)
	}

	// Print JSON result
	fmt.Println("\nRegistered Processor Result:")
	if procInfo, ok := regResult.ProcessingInfo["keyword"]; ok {
		jsonData, _ := json.MarshalIndent(procInfo, "", "  ")
		fmt.Println(string(jsonData))
	}
}
