package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Create a temporary TextItem struct for the example
type TextItem struct {
	Content string
}

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

	// Pass the processor itself as the implementations for required interfaces
	base := processor.NewBaseProcessor("keyword", provider, options, nil, p, p)
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
func (p *KeywordProcessor) HandleResponse(_ context.Context, text string, responseData interface{}) (*processor.Result, error) {
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
	keywordResult := KeywordResult{
		Keywords:    keywords,
		Categories:  categories,
		Frequencies: frequencies,
	}

	// Simple text processing: join keywords with commas
	processedText := "Keywords: " + strings.Join(keywords, ", ")

	// Create and return the result
	return &processor.Result{
		Original:  text,
		Processed: processedText,
		Data:      keywordResult,
	}, nil
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

	// Process a text
	text := "Artificial intelligence (AI) is intelligence demonstrated by machines, " +
		"as opposed to intelligence displayed by humans or other animals. Example tasks in which " +
		"AI is applied include speech recognition, computer vision, translation between natural " +
		"languages, and other mappings of inputs. AI applications include advanced web search engines, " +
		"recommendation systems, understanding human speech, self-driving cars, automated decision-making, " +
		"and competing at the highest level in strategic game systems."

	result, err := keywordProcessor.Process(context.Background(), text)
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print the result
	keywordResult, ok := result.Data.(KeywordResult)
	if !ok {
		log.Fatalf("Invalid result type")
	}

	fmt.Println("Custom Processor Result:")
	fmt.Printf("Keywords: %v\n", keywordResult.Keywords)
	fmt.Printf("Categories: %v\n", keywordResult.Categories)
	fmt.Println("Frequencies:")
	for k, v := range keywordResult.Frequencies {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// Verify that our processor was registered correctly
	regProcessor, err := processor.GetProcessor("keyword", provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor from registry: %v", err)
	}

	// Use the registered processor
	regResult, err := regProcessor.Process(context.Background(), "Machine learning is a subset of AI focused on training models to improve with experience.")
	if err != nil {
		log.Fatalf("Processing with registered processor failed: %v", err)
	}

	// Print JSON result
	fmt.Println("\nRegistered Processor Result:")
	jsonData, _ := json.MarshalIndent(regResult.Data, "", "  ")
	fmt.Println(string(jsonData))
}
