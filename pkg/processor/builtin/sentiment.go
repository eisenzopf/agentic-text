package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// SentimentResult contains the sentiment analysis results
type SentimentResult struct {
	// Sentiment is the overall sentiment (positive, negative, neutral)
	Sentiment string `json:"sentiment" default:"unknown"`
	// Score is the sentiment score (-1.0 to 1.0)
	Score float64 `json:"score" default:"0.0"`
	// Confidence is the confidence level (0.0 to 1.0)
	Confidence float64 `json:"confidence" default:"0.0"`
	// Keywords are key sentiment words from the text
	Keywords []string `json:"keywords,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("sentiment").
		WithStruct(&SentimentResult{}).
		WithRole("You are an expert sentiment analysis tool that ONLY outputs valid JSON").
		WithObjective("Analyze the sentiment expressed in the provided text accurately and objectively. Consider the overall tone, specific word choices, context, and potential nuances like sarcasm or mixed feelings").
		WithInstructions(
			"Carefully read and interpret the Input Text",
			"Determine the primary sentiment: 'positive', 'negative', or 'neutral'",
			"Assign a precise sentiment score between -1.0 (most negative) and 1.0 (most positive)",
			"Assess your confidence in the analysis on a scale of 0.0 to 1.0",
			"Extract up to 5 keywords or short phrases most representative of the sentiment",
			"Format your entire output as a single, valid JSON object conforming to the structure below",
		).
		Register()
}
