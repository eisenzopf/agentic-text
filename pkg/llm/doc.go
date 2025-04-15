/*
Package llm provides a unified interface for interacting with various Large Language Model providers.

The llm package abstracts the complexities of different LLM APIs behind a consistent interface, allowing
for easy switching between providers while maintaining a consistent codebase.

Core components:

1. Provider Interface (provider.go):
  - Provider: Main interface for interacting with LLM services
  - Generate: For generating text responses
  - GenerateJSON: For structured data generation

2. Provider Types:
  - Google (google.go): Implementation for Google's Gemini models
  - OpenAI (openai.go): Implementation for OpenAI's GPT models
  - Groq (groq.go): Implementation for Groq's models
  - Amazon (amazon.go): Implementation for Amazon Bedrock

3. Configuration:
  - Config: Standardized configuration for all providers
  - ProviderType: Enum of supported providers

4. Utilities:
  - ExtractJSONResponse: Handling JSON responses from LLMs
  - WrapWithDebugInfo: Adding debug information to responses

To use an LLM provider, create it with the appropriate configuration and use
the Provider interface methods to interact with it.
*/
package llm
