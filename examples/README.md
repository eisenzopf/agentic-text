# Agentic-Text Examples

This directory contains examples of how to use the Agentic-Text library:

## Easy Usage

The [easy_usage](./easy_usage) example demonstrates how to:
- Use the simplified `easy` package interface
- Process text with various processors via command line
- Use debug mode to see detailed model interactions
- Process multiple inputs in batch mode

Run it with:
```bash
cd easy_usage
go run main.go sentiment "I love this product"
```

## Basic Usage

The [basic_usage](./basic_usage) example demonstrates how to:
- Initialize an LLM provider
- Create a processor
- Process a single text
- Process multiple texts in parallel

Run it with:
```bash
cd basic_usage
go run main.go -processor=intent "I'd like to get my checking account balance"
```

## ProcessItem Usage

The [processitem_usage](./processitem_usage) example shows how to:
- Work with the ProcessItem data structure
- Process items individually or in batch mode
- Track processing history and metadata
- Configure model parameters via command line flags

Run it with:
```bash
cd processitem_usage
go run main.go -batch text1 "text2 with spaces" text3
```

## Custom Processor

The [custom_processor](./custom_processor) example shows how to:
- Create your own custom processor
- Register it with the processor registry
- Process text with your custom processor

Run it with:
```bash
cd custom_processor
go run main.go
```

## API Deployment

The [api_deployment](./api_deployment) example demonstrates how to:
- Create a simple REST API for text processing
- Expose multiple processors via HTTP endpoints
- Process requests in a web server

Run it with:
```bash
cd api_deployment
go run main.go
```

Test with:
```bash
# List processors
curl -X GET http://localhost:8080/api/processors

# Process text
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d '{"text": "I really enjoyed this product!", "processor": "sentiment"}'
```

## Conversation Analyzer

The [conversation](./conversation) example shows how to:
- Analyze customer service conversations from the banking domain
- Process text with multiple processors (sentiment and intent)
- Extract both emotional tone and customer intent from conversations

Run it with:
```bash
cd conversation
go run main.go
```

Or analyze your own conversation:
```bash
cd conversation
go run main.go "Your conversation text here"
``` 