package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Config represents the structure of the config.json file
type Config struct {
	Provider     string  `json:"provider"`
	Model        string  `json:"model"`
	APIKeyEnvVar string  `json:"api_key_env_var"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float64 `json:"temperature"`
}

func main() {
	// Define command-line flags
	processorType := flag.String("processor", "sentiment", "The type of processor to use (sentiment, etc.)")
	batchMode := flag.Bool("batch", false, "Process multiple text inputs as a batch")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	verbose := flag.Bool("verbose", false, "Show LLM input and output for debugging")

	// Config overrides
	providerFlag := flag.String("provider", "", "Override the LLM provider in config.json")
	modelFlag := flag.String("model", "", "Override the model name in config.json")
	apiKeyEnvFlag := flag.String("api-key-env", "", "Override the API key environment variable name in config.json")
	maxTokensFlag := flag.Int("max-tokens", 0, "Override the max tokens in config.json (0 means use config value)")
	temperatureFlag := flag.Float64("temperature", -1, "Override the temperature in config.json (-1 means use config value)")

	flag.Parse()

	// Get the text input from remaining arguments
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Error: No text provided. Please provide text to analyze.")
	}

	// Print received arguments for debugging
	fmt.Println("Received arguments:", args)

	// Join all arguments into a single text if not in batch mode
	var textToProcess string
	if !*batchMode && len(args) > 1 {
		// If there are multiple arguments and not in batch mode,
		// join them as a single text with spaces
		textToProcess = strings.Join(args, " ")
		// Replace the first arg with the joined text and clear the rest
		args = []string{textToProcess}
	}

	// Load the configuration
	configData, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	// Apply command-line overrides if provided
	if *providerFlag != "" {
		config.Provider = *providerFlag
	}
	if *modelFlag != "" {
		config.Model = *modelFlag
	}
	if *apiKeyEnvFlag != "" {
		config.APIKeyEnvVar = *apiKeyEnvFlag
	}
	if *maxTokensFlag > 0 {
		config.MaxTokens = *maxTokensFlag
	}
	if *temperatureFlag >= 0 {
		config.Temperature = *temperatureFlag
	}

	// Get API key from environment variable
	apiKey := os.Getenv(config.APIKeyEnvVar)
	if apiKey == "" {
		log.Fatalf("Error: Environment variable %s not set or empty", config.APIKeyEnvVar)
	}

	// Create provider options to capture debug info if verbose is enabled
	providerOptions := map[string]interface{}{}
	if *verbose {
		providerOptions["debug"] = true
	}

	// Update provider config with options
	providerConfig := llm.Config{
		APIKey:      apiKey,
		Model:       config.Model,
		MaxTokens:   config.MaxTokens,
		Temperature: config.Temperature,
		Options:     providerOptions,
	}

	// Map string provider to provider type
	var providerType llm.ProviderType
	switch strings.ToLower(config.Provider) {
	case "google":
		providerType = llm.Google
	case "openai":
		providerType = llm.OpenAI
	// Add more providers as they become available in the library
	// case "anthropic":
	// 	providerType = llm.Anthropic
	// case "amazon":
	// 	providerType = llm.Amazon
	// case "groq":
	// 	providerType = llm.Groq
	default:
		log.Fatalf("Unsupported provider: %s", config.Provider)
	}

	provider, err := llm.NewProvider(providerType, providerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Create a processor
	proc, err := processor.GetProcessor(*processorType, provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor: %v", err)
	}

	// Only show the prompt template if verbose is enabled
	if *verbose {
		// Create a context with our processor name for debugging
		dummyText := "SAMPLE_TEXT"
		ctx := context.Background()

		// Try to get the template by calling GeneratePrompt
		if sentimentProc, ok := proc.(*processor.SentimentProcessor); ok {
			prompt, _ := sentimentProc.GeneratePrompt(ctx, dummyText)
			fmt.Println("\n=== PROMPT TEMPLATE ===")
			fmt.Println(strings.Replace(prompt, dummyText, "<INPUT_TEXT>", 1))
			fmt.Println("=== END PROMPT TEMPLATE ===")
		} else {
			fmt.Println("Note: Prompt template is not available for this processor type")
		}
	}

	if *batchMode {
		// Process all inputs as a batch
		source := data.NewStringsSource(args)
		results, err := proc.ProcessSource(context.Background(), source, 2, 2)
		if err != nil {
			log.Fatalf("Batch processing failed: %v", err)
		}

		// Print batch results
		fmt.Println("\nResults:")
		for i, result := range results {
			// Get the original text
			origText := ""
			if item, ok := result.Original.(*data.TextItem); ok {
				origText = item.Content
			} else if s, ok := result.Original.(string); ok {
				origText = s
			}

			fmt.Printf("\nText %d: %s\n", i+1, origText)

			// Print debug info if verbose mode is enabled
			if *verbose && result.Data != nil {
				if debugInfo, ok := result.Data.(map[string]interface{})["debug"]; ok {
					if debug, ok := debugInfo.(map[string]interface{}); ok {
						// Show only the prompt and raw response
						fmt.Println("\n=== LLM INPUT ===")
						if prompt, ok := debug["prompt"].(string); ok {
							fmt.Println(prompt)
						}
						fmt.Println("=== END LLM INPUT ===")
						fmt.Println("=== LLM OUTPUT ===")
						if rawResponse, ok := debug["raw_response"].(string); ok {
							fmt.Println(rawResponse)
						}
						fmt.Println("=== END LLM OUTPUT ===")
					}
				}
			}

			// Pretty print the data
			outputData := result.Data
			// Remove debug info from output if it was already shown
			if *verbose {
				if resultMap, ok := outputData.(map[string]interface{}); ok {
					// Create a copy without the debug field
					cleanData := make(map[string]interface{})
					for k, v := range resultMap {
						if k != "debug" {
							cleanData[k] = v
						}
					}
					outputData = cleanData
				}
			}

			jsonData, _ := json.MarshalIndent(outputData, "", "  ")
			fmt.Println(string(jsonData))
		}
	} else {
		// Process just the first input
		text := args[0]
		result, err := proc.Process(context.Background(), text)
		if err != nil {
			log.Fatalf("Processing failed: %v", err)
		}

		// Print debug info if verbose mode is enabled
		if *verbose {
			if debugData, ok := result.Data.(map[string]interface{}); ok && debugData["debug"] != nil {
				if debug, ok := debugData["debug"].(map[string]interface{}); ok {
					// Show only the prompt and raw response
					fmt.Println("\n=== LLM INPUT ===")
					if prompt, ok := debug["prompt"].(string); ok {
						fmt.Println(prompt)
					}
					fmt.Println("=== END LLM INPUT ===")

					fmt.Println("=== LLM OUTPUT ===")
					if rawResponse, ok := debug["raw_response"].(string); ok {
						fmt.Println(rawResponse)
					}
					fmt.Println("=== END LLM OUTPUT ===")
				}
			} else {
				fmt.Println("No debug information available")
			}
		}

		// Print the result as JSON (final processed result)
		outputData := result.Data
		// Remove debug info from output
		if *verbose {
			if resultMap, ok := outputData.(map[string]interface{}); ok {
				// Create a copy without the debug field
				cleanData := make(map[string]interface{})
				for k, v := range resultMap {
					if k != "debug" {
						cleanData[k] = v
					}
				}
				outputData = cleanData
			}
		}

		jsonData, _ := json.MarshalIndent(outputData, "", "  ")
		fmt.Println(string(jsonData))
	}
}
