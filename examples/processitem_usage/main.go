package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
	"github.com/eisenzopf/agentic-text/pkg/pipeline"
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
	// Explicitly register processors for attributes analysis
	processor.Register("required_attributes", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
		return processor.NewRequiredAttributesProcessor(provider, options)
	})

	processor.Register("get_attributes", func(provider llm.Provider, options processor.Options) (processor.Processor, error) {
		return processor.NewAttributeProcessor(provider, options)
	})

	// Define command-line flags
	processorType := flag.String("processor", "required_attributes", "The type of processor to use (required_attributes, get_attributes, etc.)")
	secondaryProcessorType := flag.String("secondary", "", "Optional secondary processor to chain after the primary")
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
	default:
		log.Fatalf("Unsupported provider: %s", config.Provider)
	}

	provider, err := llm.NewProvider(providerType, providerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}

	// Create primary processor
	primaryProc, err := processor.Create(*processorType, provider, processor.Options{})
	if err != nil {
		log.Fatalf("Failed to get processor: %v", err)
	}

	// Create a processor chain if secondary processor is specified
	var processors []processor.Processor
	processors = append(processors, primaryProc)

	if *secondaryProcessorType != "" {
		secondaryProc, err := processor.Create(*secondaryProcessorType, provider, processor.Options{})
		if err != nil {
			log.Fatalf("Failed to get secondary processor: %v", err)
		}
		processors = append(processors, secondaryProc)
		fmt.Printf("Using processor chain: %s -> %s\n", *processorType, *secondaryProcessorType)
	} else {
		fmt.Printf("Using processor: %s\n", *processorType)
	}

	// Create a chain
	// When chaining required_attributes -> get_attributes:
	// 1. required_attributes analyzes the text and determines what data attributes would be needed
	// 2. required_attributes returns a JSON with an array of attribute definitions
	// 3. get_attributes receives this JSON input, extracts the attribute definitions
	// 4. get_attributes then tries to extract values for those attributes from the original text
	// This demonstrates how processors can work together in a pipeline to build on each other's output
	chain := pipeline.NewChain("analysis-chain", processors...)

	if *batchMode {
		// Create ProcessItems for all inputs and process them
		items := make([]*data.ProcessItem, len(args))
		for i, text := range args {
			items[i] = data.NewTextProcessItem(fmt.Sprintf("text-%d", i+1), text, nil)
		}

		// Create a ProcessItemSource from the items
		source := data.NewProcessItemSliceSource(items)

		// Process the batch
		results, err := chain.ProcessSource(context.Background(), source, 2, 2)
		if err != nil {
			log.Fatalf("Batch processing failed: %v", err)
		}

		// Print batch results
		fmt.Println("\nResults:")
		for i, result := range results {
			fmt.Printf("\nItem %d (ID: %s):\n", i+1, result.ID)
			fmt.Printf("Content Type: %s\n", result.ContentType)

			// Print the content based on type
			switch result.ContentType {
			case "text":
				text, _ := result.GetTextContent()
				fmt.Printf("Content: %s\n", text)
			case "json":
				jsonData, _ := json.MarshalIndent(result.Content, "", "  ")
				fmt.Printf("Content: %s\n", jsonData)
			default:
				fmt.Printf("Content: %v\n", result.Content)
			}

			// Print processing info
			fmt.Println("Processing History:")
			for procName, procInfo := range result.ProcessingInfo {
				fmt.Printf("  Processor: %s\n", procName)

				// Print debug info if verbose mode is enabled
				if *verbose {
					if debugMap, ok := procInfo.(map[string]interface{}); ok {
						if debug, ok := debugMap["debug"].(map[string]interface{}); ok {
							fmt.Printf("    --- LLM INPUT ---\n")
							if prompt, ok := debug["prompt"].(string); ok {
								fmt.Printf("    %s\n", prompt)
							}
							fmt.Printf("    --- LLM OUTPUT ---\n")
							if rawResponse, ok := debug["raw_response"].(string); ok {
								fmt.Printf("    %s\n", rawResponse)
							}
						}
					}
				}

				// Create a clean copy of the processing info without debug data
				var cleanInfo interface{} = procInfo
				if *verbose {
					if infoMap, ok := procInfo.(map[string]interface{}); ok {
						cleanInfoMap := make(map[string]interface{})
						for k, v := range infoMap {
							if k != "debug" {
								cleanInfoMap[k] = v
							}
						}
						cleanInfo = cleanInfoMap
					}
				}

				// Print processing info
				infoJSON, _ := json.MarshalIndent(cleanInfo, "    ", "  ")
				fmt.Printf("    %s\n", infoJSON)
			}
		}
	} else {
		// Process just the first input
		text := args[0]

		// Create a ProcessItem directly
		item := data.NewTextProcessItem("input-1", text, map[string]interface{}{
			"source":    "command-line",
			"timestamp": fmt.Sprintf("%v", float64(float32(time.Now().Unix()))),
		})

		// Process the item
		result, err := chain.Process(context.Background(), item)
		if err != nil {
			log.Fatalf("Processing failed: %v", err)
		}

		// Print the result
		fmt.Println("\nResult:")
		fmt.Printf("Content Type: %s\n", result.ContentType)

		// Print the content based on type
		switch result.ContentType {
		case "text":
			text, _ := result.GetTextContent()
			fmt.Printf("Content: %s\n", text)
		case "json":
			jsonData, _ := json.MarshalIndent(result.Content, "", "  ")
			fmt.Printf("Content: %s\n", jsonData)
		default:
			fmt.Printf("Content: %v\n", result.Content)
		}

		// Print processing info
		fmt.Println("Processing History:")
		for procName, procInfo := range result.ProcessingInfo {
			fmt.Printf("  Processor: %s\n", procName)

			// Print debug info if verbose mode is enabled
			if *verbose {
				if debugMap, ok := procInfo.(map[string]interface{}); ok {
					if debug, ok := debugMap["debug"].(map[string]interface{}); ok {
						fmt.Printf("    --- LLM INPUT ---\n")
						if prompt, ok := debug["prompt"].(string); ok {
							fmt.Printf("    %s\n", prompt)
						}
						fmt.Printf("    --- LLM OUTPUT ---\n")
						if rawResponse, ok := debug["raw_response"].(string); ok {
							fmt.Printf("    %s\n", rawResponse)
						}
					}
				}
			}

			// Create a clean copy of the processing info without debug data
			var cleanInfo interface{} = procInfo
			if *verbose {
				if infoMap, ok := procInfo.(map[string]interface{}); ok {
					cleanInfoMap := make(map[string]interface{})
					for k, v := range infoMap {
						if k != "debug" {
							cleanInfoMap[k] = v
						}
					}
					cleanInfo = cleanInfoMap
				}
			}

			// Print processing info
			infoJSON, _ := json.MarshalIndent(cleanInfo, "    ", "  ")
			fmt.Printf("    %s\n", infoJSON)
		}
	}
}
