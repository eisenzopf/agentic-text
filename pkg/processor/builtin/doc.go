// Package builtin provides a collection of ready-to-use processors
// for common text analysis tasks using the processor framework.
//
// These processors are registered automatically when the package is imported,
// so you can immediately use them with processor.Create()
//
// Available processors:
// - sentiment: Analyzes the sentiment of text, returning sentiment type, score, confidence, and keywords
// - intent: Identifies the primary intent in customer service conversations
// - keyword_extraction: Extracts important keywords from text with relevance scores and categories
// - required_attributes: Identifies data attributes required to answer a set of questions
// - get_attributes: Extracts attribute values from text based on the identified attributes
package builtin
