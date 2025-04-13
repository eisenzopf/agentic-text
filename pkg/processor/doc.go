// Package processor provides a framework for implementing text processors using LLMs.
//
// This package contains the core interfaces, types, and utilities needed to create
// processors that analyze text using large language models.
//
// The framework handles:
// - Processing text through LLMs with customizable prompts
// - Parsing and validating LLM responses
// - Converting responses to structured data
// - Managing processor registration and discovery
//
// To implement your own processor:
// 1. Define your result struct (e.g., MyResult)
// 2. Create a prompt generator that implements PromptGenerator
// 3. Register your processor with RegisterGenericProcessor
//
// Ready-to-use processors are available in the builtin subpackage.
// To use these builtin processors, import:
//
//	"github.com/eisenzopf/agentic-text/pkg/processor/builtin"
package processor
