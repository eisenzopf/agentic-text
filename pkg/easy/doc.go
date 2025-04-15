/*
Package easy provides simplified, one-liner functions for common text processing operations.

The easy package is designed for users who want a quick and simple way to use the agentic-text
functionality without dealing with the complexity of the underlying processor framework.

Core components:

1. Configuration (easy.go):
  - Config: Simplified configuration structure
  - DefaultConfig: Sensible default configuration values

2. ProcessorWrapper (easy.go):
  - ProcessorWrapper: Handles the creation and management of processors
  - Process: For processing single text items
  - ProcessBatch: For processing multiple text items in parallel

3. Convenience Functions (utils.go):
  - Sentiment: One-liner for sentiment analysis
  - Intent: One-liner for intent detection
  - ProcessText: Generic text processing
  - ProcessBatchText: Batch processing of multiple texts
  - PrettyPrint: For formatting results as JSON

This package abstracts away the creation of providers, processors, and data structures,
making it ideal for simple applications or quick prototyping.
*/
package easy
