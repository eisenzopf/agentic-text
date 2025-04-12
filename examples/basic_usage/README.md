# Basic Usage Example

This example demonstrates the basic programming interface for the agentic-text library, showing how to initialize an LLM provider and create processors.

## Features

- Initialize an LLM provider
- Create a processor
- Process a single text
- Process multiple texts in parallel
- Configure processing via command line flags

## Usage

From within the directory, run:

```bash
go run main.go -processor=intent "I'd like to get my checking account balance"
```

## Command Line Options

The example supports several command line flags:

```bash
# Change the processor type
go run main.go -processor=sentiment "I love this product"

# Process multiple texts
go run main.go -processor=intent "Check my balance" "Transfer funds" "Report fraud"

# Use a custom configuration file
go run main.go -config=./custom_config.json -processor=intent "Check my balance"
```

## Configuration

The example uses a `config.json` file to configure the LLM provider. You can modify this file to change:

- Model provider (Google, OpenAI, etc.)
- API keys
- Model parameters
- Default processor settings

## Code Structure

The main.go file shows:
1. How to initialize the LLM service
2. How to create different types of processors
3. Proper error handling
4. Processing single and multiple inputs 