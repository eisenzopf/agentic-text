package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/easy"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

func main() {
	// Default values
	processorType := "sentiment"
	debugMode := false

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) == 0 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {
		printUsage()
		return
	}

	// Extract flags and collect non-flag arguments
	var inputs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-debug" || args[i] == "--debug" {
			debugMode = true
		} else if args[i] == "-h" || args[i] == "--help" {
			printUsage()
			return
		} else if strings.HasPrefix(args[i], "-") {
			fmt.Printf("Unknown flag: %s\n", args[i])
			printUsage()
			return
		} else if i == 0 {
			// First non-flag argument is the processor type
			processorType = args[i]
		} else {
			// All other non-flag arguments are inputs
			inputs = append(inputs, args[i])
		}
	}

	// Check if processor exists
	availableProcessors := easy.ListAvailableProcessors()
	processorExists := false
	for _, p := range availableProcessors {
		if p == processorType {
			processorExists = true
			break
		}
	}

	if !processorExists {
		fmt.Printf("Error: Processor '%s' not found.\n", processorType)
		fmt.Printf("Available processors: %v\n", availableProcessors)
		os.Exit(1)
	}

	// Create configuration if debug mode is enabled
	var config *easy.Config
	if debugMode {
		config = &easy.Config{
			Provider:    llm.Google,
			Model:       "gemini-2.0-flash",
			MaxTokens:   1024,
			Temperature: 0.2,
			Debug:       true,
		}
	}

	if len(inputs) == 0 {
		// No input provided, use default
		inputs = append(inputs, "I absolutely love this product")
		fmt.Printf("No input provided. Using default: '%s'\n", inputs[0])
	}

	fmt.Printf("Available processors: %v\n", availableProcessors)
	fmt.Printf("Using processor: '%s'\n", processorType)

	if len(inputs) == 1 {
		// Single input mode
		fmt.Printf("Processing single input: '%s'\n\n", inputs[0])

		var result map[string]interface{}
		var err error

		if debugMode {
			result, err = easy.ProcessTextWithConfig(inputs[0], processorType, config)
		} else {
			result, err = easy.ProcessText(inputs[0], processorType)
		}

		if err != nil {
			log.Fatalf("Processing failed: %v", err)
		}

		prettyResult, err := easy.PrettyPrint(result)
		if err != nil {
			log.Fatalf("Failed to format result: %v", err)
		}
		fmt.Printf("Result:\n%s\n", prettyResult)
	} else {
		// Batch mode
		fmt.Printf("Batch processing %d inputs\n\n", len(inputs))

		var results []map[string]interface{}
		var err error

		if debugMode {
			// Create a wrapper with config for batch processing
			wrapper, err := easy.NewWithConfig(processorType, config)
			if err != nil {
				log.Fatalf("Failed to create processor: %v", err)
			}

			results, err = wrapper.ProcessBatch(inputs, 2) // Use concurrency of 2
		} else {
			results, err = easy.ProcessBatchText(inputs, processorType, 2) // Use concurrency of 2
		}

		if err != nil {
			log.Fatalf("Batch processing failed: %v", err)
		}

		for i, result := range results {
			prettyResult, err := easy.PrettyPrint(result)
			if err != nil {
				log.Fatalf("Failed to format result for input %d: %v", i+1, err)
			}
			fmt.Printf("Result for input %d: '%s'\n%s\n\n", i+1, inputs[i], prettyResult)
		}
	}
}

func printUsage() {
	// Get the program name - use a clean base name instead of full path
	progName := "easy_usage"
	if len(os.Args) > 0 {
		// Extract just the base filename without path
		progName = filepath.Base(os.Args[0])
		// If it's a temporary go run path, simplify it
		if strings.Contains(progName, "go-build") {
			progName = "easy_usage"
		}
	}

	fmt.Printf(`Usage: %s [options] processor_type [input1] [input2] ...

Options:
  -debug, --debug    Enable debug mode to see prompts and raw responses
  -h, --help         Show this help message

Processor Types: %v

Examples:
  %s sentiment "I love this product"
  %s intent "I want to cancel my subscription"
  %s sentiment -debug "I love this product"
  
  # Batch processing multiple inputs
  %s sentiment "I love this product" "This product is terrible" "It is okay I guess"
  
If no input is provided, it defaults to "I absolutely love this product"
`, progName, easy.ListAvailableProcessors(), progName, progName, progName, progName)
}
