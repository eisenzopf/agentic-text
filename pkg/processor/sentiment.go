package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// SentimentProcessor analyzes the sentiment of text
type SentimentProcessor struct {
	// Embed BaseProcessor to inherit all methods
	BaseProcessor
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
	p := &SentimentProcessor{}

	// Create client from provider
	client := llm.NewProviderClient(provider)

	// Create and embed base processor
	base := NewBaseProcessor("sentiment", []string{"text"}, client, nil, p, p, options)
	p.BaseProcessor = *base

	return p, nil
}

// GeneratePrompt implements PromptGenerator interface - generates the sentiment analysis prompt
func (p *SentimentProcessor) GeneratePrompt(ctx context.Context, text string) (string, error) {
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
{
  "sentiment": "...", // "positive", "negative", or "neutral"
  "score": ...,     // Float between -1.0 and 1.0
  "confidence": ..., // Float between 0.0 and 1.0
  "keywords": ["...", "..."] // Array of up to 5 strings
}`, text), nil
}

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *SentimentProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	// Check if responseData is a string, which can happen with some providers
	if strResponse, ok := responseData.(string); ok {
		// Remove markdown code block if present
		cleanResponse := strResponse
		if strings.HasPrefix(cleanResponse, "```json") && strings.HasSuffix(cleanResponse, "```") {
			// Extract content between ```json and ```
			cleanResponse = strings.TrimPrefix(cleanResponse, "```json")
			cleanResponse = strings.TrimSuffix(cleanResponse, "```")
			cleanResponse = strings.TrimSpace(cleanResponse)
		} else if strings.HasPrefix(cleanResponse, "```") && strings.HasSuffix(cleanResponse, "```") {
			// Extract content between ``` and ```
			cleanResponse = strings.TrimPrefix(cleanResponse, "```")
			cleanResponse = strings.TrimSuffix(cleanResponse, "```")
			cleanResponse = strings.TrimSpace(cleanResponse)
		}

		// Try to parse the string as JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(cleanResponse), &data); err != nil {
			// If parsing fails, wrap it as a response
			return map[string]interface{}{
				"sentiment":      "unknown",
				"score":          0.0,
				"confidence":     0.0,
				"keywords":       []string{},
				"response":       strResponse,
				"processor_type": "sentiment",
			}, nil
		}
		// If parsing succeeds, use the parsed data
		responseData = data
	}

	// Convert the response data to map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format: %T", responseData)
	}

	// Check if debug info exists and preserve it
	var debugInfo interface{}
	if debug, exists := data["debug"]; exists {
		debugInfo = debug
	}

	// Check if we got a non-JSON response wrapped in a "response" field
	if response, exists := data["response"]; exists && len(data) <= 2 { // data has only response and maybe debug
		// This is a fallback case where the LLM didn't produce valid JSON
		// Create a placeholder sentiment result
		resultMap := map[string]interface{}{
			"sentiment":      "unknown",
			"score":          0.0,
			"confidence":     0.0,
			"keywords":       []string{},
			"response":       response,
			"processor_type": "sentiment",
		}

		// Add debug info back if it existed
		if debugInfo != nil {
			resultMap["debug"] = debugInfo
		}

		return resultMap, nil
	}

	// Normal case - we have sentiment data fields
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

	// Create result map with sentiment data
	resultMap := map[string]interface{}{
		"sentiment":      sentiment,
		"score":          score,
		"confidence":     confidence,
		"keywords":       keywords,
		"processor_type": "sentiment",
	}

	// Add debug info back if it existed
	if debugInfo != nil {
		resultMap["debug"] = debugInfo
	}

	return resultMap, nil
}

// Register the processor with the registry
func init() {
	Register("sentiment", func(provider llm.Provider, options Options) (Processor, error) {
		return NewSentimentProcessor(provider, options)
	})
}
