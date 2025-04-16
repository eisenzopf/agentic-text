# API Deployment Example

This example demonstrates how to create a simple REST API for the agentic-text library, allowing you to expose text processing capabilities via HTTP endpoints.

## Features

- Create a simple REST API for text processing
- Expose multiple processors via HTTP endpoints
- Process requests in a web server
- JSON response formatting
- Error handling for web requests

## Setup

Before running the server, set your Gemini API key in the environment:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

## Usage

From within the directory, run:

```bash
go run main.go
```

This will start a web server on port 8080 (or the port specified in the `PORT` environment variable). You can then use the API with curl, Postman, or any HTTP client.

## API Endpoints

### List Available Processors

```bash
# List all available processors
curl -X GET http://localhost:8080/api/processors
```

Response:
```json
{
  "processors": ["sentiment", "intent", "keyword_extraction", "speech_act", "get_attributes", "required_attributes"],
  "count": 6
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
  "original": "I really enjoyed this product!",
  "result": {
    "sentiment": "positive",
    "score": 0.8,
    "confidence": 0.95,
    "keywords": ["enjoyed", "really"],
    "processor_type": "sentiment"
  },
  "success": true
}
```

## Example Script

The repository includes a `process_examples.sh` script that demonstrates how to use all available processors with sample inputs:

```bash
# Make the script executable
chmod +x process_examples.sh

# Run all examples
./process_examples.sh
```

The script runs examples for:
- Sentiment analysis
- Intent detection
- Keyword extraction
- Speech act classification
- Attribute extraction
- Required attributes identification

## Configuration

The API server configuration can be modified in the main.go file:
- Port number (via `PORT` environment variable)
- LLM provider settings (model, max tokens, temperature)
- API key (via `GEMINI_API_KEY` environment variable)

## Production Considerations

For production use, consider:
- Adding authentication
- Implementing rate limiting
- Setting up HTTPS
- Adding monitoring and logging
- Deploying behind a reverse proxy 