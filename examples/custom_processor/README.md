# Custom Processor Example

This example demonstrates how to create your own custom processor with the agentic-text library to implement specialized text processing functionality.

## Features

- Create your own custom processor
- Register it with the processor registry
- Process text with your custom processor
- Implement custom logic and output formats
- Handle processor-specific configuration

## Usage

From within the directory, run:

```bash
go run main.go
```

## How It Works

The example:
1. Defines a custom processor structure
2. Implements the required processor interface methods
3. Registers the processor with the global registry
4. Processes sample text using the custom processor

## Creating Custom Processors

The code shows how to:
- Define a processor struct that implements the Processor interface
- Implement key methods like Process() and Name()
- Add custom configuration
- Register your processor so it can be used throughout the application
- Handle processor-specific logic, prompts, and result parsing

## Custom Prompt Engineering

The example demonstrates how to create effective prompts for your LLM, including:
- Structuring system and user prompts
- Specifying desired output formats
- Adding context and examples
- Handling different types of inputs

## Extending The Example

You can modify this example to create:
- Domain-specific analyzers (legal, medical, financial)
- Custom classification processors
- Entity extraction processors
- Specialized formatting processors
- Multi-step processing pipelines 