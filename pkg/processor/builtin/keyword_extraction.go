package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Keyword represents a single extracted keyword
type Keyword struct {
	// Term is the extracted keyword term
	Term string `json:"term"`
	// Relevance is the relevance score from 0.0 to 1.0
	Relevance float64 `json:"relevance"`
	// Category is the category of the keyword (e.g., "topic", "person", "location")
	Category string `json:"category"`
}

// KeywordResult contains the keyword extraction results
type KeywordResult struct {
	// Keywords is an array of extracted keywords
	Keywords []Keyword `json:"keywords,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("keyword_extraction").
		WithStruct(&KeywordResult{}).
		WithRole("You are an expert at extracting important keywords from text").
		WithObjective("Analyze the provided text and extract the most meaningful keywords").
		WithInstructions(
			"Carefully read and interpret the Input Text",
			"Extract the most important keywords or key phrases that represent the main topics",
			"For each keyword, provide the keyword term, relevance score (0.0 to 1.0), and category",
			"Categories include: 'topic', 'person', 'location', 'concept', 'organization'",
			"Format your entire output as a single, valid JSON object",
		).
		WithValidation().
		Register()
}
