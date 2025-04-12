# API Deployment Example

This example demonstrates how to create a simple REST API for the agentic-text library, allowing you to expose text processing capabilities via HTTP endpoints.

## Features

- Create a simple REST API for text processing
- Expose multiple processors via HTTP endpoints
- Process requests in a web server
- JSON response formatting
- Error handling for web requests

## Usage

From within the directory, run:

```bash
go run main.go
```

This will start a web server on port 8080. You can then use the API with curl, Postman, or any HTTP client.

## API Endpoints

### List Available Processors

```bash
# List all available processors
curl -X GET http://localhost:8080/api/processors
```

Response:
```json
{
  "processors": ["sentiment", "intent", "summarize", "translate", "topic"]
}
```

### Process Text

```bash
# Process text with a specific processor
curl -X POST http://localhost:8080/api/process \
  -H "Content-Type: application/json" \
  -d '{"text": "I really enjoyed this product!", "processor": "sentiment"}'
```

Response:
```json
{
  "input": "I really enjoyed this product!",
  "processor": "sentiment",
  "result": {
    "sentiment": "positive",
    "confidence": 0.92,
    "score": 0.85
  },
  "processing_time_ms": 245
}
```

## Configuration

The API server configuration can be modified in the main.go file:
- Port number
- Allowed processors
- Request timeouts
- CORS settings

## Production Considerations

For production use, consider:
- Adding authentication
- Implementing rate limiting
- Setting up HTTPS
- Adding monitoring and logging
- Deploying behind a reverse proxy 