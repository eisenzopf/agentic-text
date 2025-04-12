# ProcessItem Usage Example

This example demonstrates how to work with the `ProcessItem` data structure in the agentic-text library, which provides additional metadata and history tracking for processed items.

## Features

- Work with the ProcessItem data structure
- Process items individually or in batch mode
- Track processing history and metadata
- Configure model parameters via command line flags
- Handle processing state and errors

## Usage

From within the directory, run:

```bash
go run main.go -batch text1 "text2 with spaces" text3
```

## Command Line Options

The example supports several command line flags:

```bash
# Process a single item
go run main.go "Check my balance"

# Process multiple items in batch mode
go run main.go -batch "I love this product" "This is terrible" "It's okay I guess"

# Select a specific processor
go run main.go -processor=sentiment "I love this product"

# Use a custom configuration file
go run main.go -config=./custom_config.json "Check my balance"
```

## ProcessItem Structure

The ProcessItem structure provides:
- Input text
- Processed result
- Processing metadata (timestamps, model info)
- Processing history (all steps and changes)
- Error information
- Custom metadata fields

## Advanced Usage

The example shows how to:
1. Create ProcessItem objects
2. Submit them for processing
3. Handle processing results
4. Access metadata and history
5. Process items with different configurations

## Configuration

The example uses a `config.json` file to configure the LLM provider, which you can customize for your needs. 