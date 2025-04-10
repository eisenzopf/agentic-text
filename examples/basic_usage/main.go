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

func main() {
	// Initialize LLM provider (using a mock implementation for the example)
	config := llm.Config{
		APIKey:      "your-api-key", // In a real app, get this from environment variables
		Model:       "gemini-pro",
		MaxTokens:   1024,
		Temperature: 0.2,
	}

	provider, err := llm.NewProvider(llm.Google, config)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Create a processor
	sentimentProcessor, err := processor.GetProcessor("sentiment", provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor: %v", err)
	}

	// Process a single text
	result, err := sentimentProcessor.Process(context.Background(), "I absolutely love this product! It's amazing!")
	if err != nil {
		log.Fatalf("Processing failed: %v", err)
	}

	// Print the result
	sentimentResult, ok := result.Data.(processor.SentimentResult)
	if !ok {
		log.Fatalf("Invalid result type")
	}

	fmt.Printf("Sentiment: %s\n", sentimentResult.Sentiment)
	fmt.Printf("Score: %.2f\n", sentimentResult.Score)
	fmt.Printf("Confidence: %.2f\n", sentimentResult.Confidence)
	fmt.Printf("Keywords: %v\n", sentimentResult.Keywords)

	// Process multiple texts
	texts := []string{
		"I'm really disappointed with this service.",
		"The product is okay, but nothing special.",
		"This is the best experience I've ever had!",
	}

	// Create a data source from the texts
	source := data.NewStringsSource(texts)

	// Process the source
	results, err := sentimentProcessor.ProcessSource(context.Background(), source, 2, 2)
	if err != nil {
		log.Fatalf("Batch processing failed: %v", err)
	}

	// Print batch results
	fmt.Println("\nBatch Results:")
	for i, result := range results {
		// Get the original text
		origText := ""
		if item, ok := result.Original.(*data.TextItem); ok {
			origText = item.Content
		} else if s, ok := result.Original.(string); ok {
			origText = s
		}

		fmt.Printf("\nText %d: %s\n", i+1, origText)

		// Pretty print the data
		jsonData, _ := json.MarshalIndent(result.Data, "", "  ")
		fmt.Println(string(jsonData))
	}
}
