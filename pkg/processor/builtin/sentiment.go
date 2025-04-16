package builtin

import (
	"context"
	"fmt"

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

// SentimentPrompt is a prompt generator for sentiment analysis
type SentimentPrompt struct{}

// GeneratePrompt implements PromptGenerator interface
func (p *SentimentPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	// Generate example JSON from the result struct
	exampleResult := &SentimentResult{}
	jsonExample := processor.GenerateJSONExample(exampleResult)

	return fmt.Sprintf(`**Role:** You are an expert sentiment analysis tool that ONLY outputs valid JSON.

**Objective:** Analyze the sentiment expressed in the provided text accurately and objectively. Consider the overall tone, specific word choices, context, and potential nuances like sarcasm or mixed feelings.

**Input Text:**
%s

**Instructions:**
1.  Carefully read and interpret the Input Text.
2.  Determine the primary sentiment: "positive", "negative", or "neutral".
3.  Assign a precise sentiment score between -1.0 (most negative) and 1.0 (most positive).
4.  Assess your confidence in the analysis on a scale of 0.0 to 1.0.
5.  Extract up to 5 keywords or short phrases most representative of the sentiment.
6.  Format your entire output as a single, valid JSON object conforming to the structure below.
7.  *** IMPORTANT: Your ENTIRE response must be a single JSON object, without ANY additional text, explanation, or markdown formatting. ***

**Required JSON Output Structure:**
%s`, text, jsonExample), nil
}

// Register the processor with the registry
func init() {
	// Register the sentiment processor using the generic processor registration
	processor.RegisterGenericProcessor(
		"sentiment",        // name
		[]string{"text"},   // contentTypes
		&SentimentResult{}, // resultStruct
		&SentimentPrompt{}, // promptGenerator
		nil,                // no custom initialization needed
		false,              // No struct validation needed by default
	)
}
