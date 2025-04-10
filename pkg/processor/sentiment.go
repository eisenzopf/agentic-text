package processor

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// SentimentProcessor analyzes the sentiment of text
type SentimentProcessor struct {
	*BaseProcessor
}

// SentimentResult contains the sentiment analysis results
type SentimentResult struct {
	// Sentiment is the overall sentiment (positive, negative, neutral)
	Sentiment string `json:"sentiment"`
	// Score is the sentiment score (-1.0 to 1.0)
	Score float64 `json:"score"`
	// Confidence is the confidence level (0.0 to 1.0)
	Confidence float64 `json:"confidence"`
	// Keywords are key sentiment words from the text
	Keywords []string `json:"keywords,omitempty"`
}

// NewSentimentProcessor creates a new sentiment processor
func NewSentimentProcessor(provider llm.Provider, options Options) (*SentimentProcessor, error) {
	base := NewBaseProcessor("sentiment", provider, options)
	return &SentimentProcessor{
		BaseProcessor: base,
	}, nil
}

// GeneratePrompt overrides the BaseProcessor's method
func (p *SentimentProcessor) GeneratePrompt(_ context.Context, text string) (string, error) {
	return fmt.Sprintf(`Analyze the sentiment of the following text. Be as accurate as possible.
Text: %s

Respond with a JSON object containing:
- "sentiment": The overall sentiment as "positive", "negative", or "neutral"
- "score": A sentiment score from -1.0 (very negative) to 1.0 (very positive)
- "confidence": Your confidence in this analysis from 0.0 to 1.0
- "keywords": An array of up to 5 key sentiment-expressing words from the text

Format your response as valid JSON.`, text), nil
}

// PostProcess overrides the BaseProcessor's method
func (p *SentimentProcessor) PostProcess(_ context.Context, text string, responseData interface{}) (*Result, error) {
	// Convert the response data to SentimentResult
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format")
	}

	// Extract sentiment fields
	sentiment, _ := data["sentiment"].(string)
	score, _ := data["score"].(float64)
	confidence, _ := data["confidence"].(float64)

	// Extract keywords
	var keywords []string
	if keywordsData, ok := data["keywords"].([]interface{}); ok {
		for _, k := range keywordsData {
			if keyword, ok := k.(string); ok {
				keywords = append(keywords, keyword)
			}
		}
	}

	// Create sentiment result
	sentimentResult := SentimentResult{
		Sentiment:  sentiment,
		Score:      score,
		Confidence: confidence,
		Keywords:   keywords,
	}

	return &Result{
		Original:  text,
		Processed: text,
		Data:      sentimentResult,
	}, nil
}

// Register the processor with the registry
func init() {
	Register("sentiment", func(provider llm.Provider, options Options) (Processor, error) {
		return NewSentimentProcessor(provider, options)
	})
}
