# Easy Usage Example

This example demonstrates how to use the simplified `easy` package interface of the agentic-text library to process text with various processors.

## Features

- Use the simplified `easy` package interface
- Process text with various processors via command line
- Use debug mode to see detailed model interactions
- Process multiple inputs in batch mode

## Usage

From within the directory, run:

```bash
go run main.go sentiment "I love this product"
```

You can use different processors by changing the first argument:

```bash
go run main.go intent "I'd like to get my checking account balance"
```

## Batch Mode

Process multiple inputs at once:

```bash
go run main.go sentiment "I love this product" "This is terrible" "It's okay I guess"
```

## Debug Mode

Enable debug mode to see detailed model interactions:

```bash
go run main.go -debug sentiment "I love this product"
```

## Available Processors

The example supports various processors including:
- sentiment: Analyze the emotional tone of text
- intent: Determine the user's intention
- summarize: Create a concise summary
- translate: Convert text to another language
- topic: Identify the main topic of discussion

Check the main.go file for a complete list of available processors. 