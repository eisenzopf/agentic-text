/*
Package pipeline provides tools for chaining multiple text processors together.

The pipeline package allows for creating processing pipelines where the output of one processor
becomes the input to the next, enabling complex text processing workflows with minimal code.

Core components:

1. Chain (chain.go):
  - Chain: Main structure for processor chains
  - Process: Method for processing a single item through the chain
  - ProcessBatch: Method for batch processing items through the chain
  - ProcessSource: Method for processing a data source through the chain

Using pipelines allows for modular, composable text processing workflows where each step
is handled by a specialized processor.
*/
package pipeline
