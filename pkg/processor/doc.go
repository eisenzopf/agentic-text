/*
Package processor provides a framework for creating text processing pipelines using LLMs.

The processor package allows for creating modular, composable text processors that can:
- Pre-process text before sending to LLMs
- Generate prompts for LLMs
- Parse and handle LLM responses
- Extract structured data from unstructured text

Core components:

1. Interfaces (interfaces.go):
  - Processor: Main interface for processing items
  - TextPreProcessor: For pre-processing text before LLM
  - PromptGenerator: For generating LLM prompts
  - ResponseHandler: For handling LLM responses

2. Base Processors (base_processor.go):
  - BaseProcessor: Provides core implementation of the Processor interface
  - Handles common operations like content extraction and LLM calling

3. Generic Processors (generic_processor.go):
  - GenericProcessor: Extends BaseProcessor with standard response handling
  - RegisterGenericProcessor: Helper for registering processors

4. Response Handling (response_handler.go):
  - BaseResponseHandler: Provides common response handling functionality
  - Includes JSON parsing, field mapping, and validation

5. Utilities:
  - JSON utilities (json_utils.go): Tools for working with JSON data
  - Validation (validation.go): Functions for validating LLM responses

6. Registry (registry.go):
  - Register: Registers processor factories
  - Create: Creates processors by name

To create a custom processor, implement the required interfaces and register
your processor factory using Register() or use the RegisterGenericProcessor()
helper function for common cases.
*/
package processor
