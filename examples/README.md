# Agentic-Text Examples

This directory contains examples of how to use the Agentic-Text library:

## Basic Usage

The [basic_usage](./basic_usage) example demonstrates how to:
- Initialize an LLM provider
- Create a processor
- Process a single text
- Process multiple texts in parallel

Run it with:
```bash
cd basic_usage
go run main.go
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