package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// SentimentProcessor analyzes the sentiment of text
type SentimentProcessor struct {
	base BaseProcessor
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
		base: *base,
	}, nil
}

// GetName returns the processor name
func (p *SentimentProcessor) GetName() string {
	return p.base.GetName()
}

// GeneratePrompt generates the sentiment analysis prompt
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

// PreProcess delegates to the base processor
func (p *SentimentProcessor) PreProcess(ctx context.Context, text string) (string, error) {
	return p.base.PreProcess(ctx, text)
}

// Process processes a single text item
func (p *SentimentProcessor) Process(ctx context.Context, text string) (*Result, error) {
	// Pre-process the text
	processedText, err := p.PreProcess(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("pre-processing error: %w", err)
	}

	// Generate the prompt using our specific implementation
	prompt, err := p.GeneratePrompt(ctx, processedText)
	if err != nil {
		return nil, fmt.Errorf("prompt generation error: %w", err)
	}

	// Get result from LLM
	var responseData interface{}
	err = p.base.provider.GenerateJSON(ctx, prompt, &responseData)
	if err != nil {
		return nil, fmt.Errorf("LLM error: %w", err)
	}

	// Post-process the result
	result, err := p.PostProcess(ctx, processedText, responseData)
	if err != nil {
		return nil, fmt.Errorf("post-processing error: %w", err)
	}

	// Set the original text in the result
	result.Original = text

	return result, nil
}

// PostProcess handles the LLM response
func (p *SentimentProcessor) PostProcess(ctx context.Context, text string, responseData interface{}) (*Result, error) {
	// Convert the response data to map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response data format")
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
			"sentiment":  "unknown",
			"score":      0.0,
			"confidence": 0.0,
			"keywords":   []string{},
			"response":   response,
		}

		// Add debug info back if it existed
		if debugInfo != nil {
			resultMap["debug"] = debugInfo
		}

		return &Result{
			Original:  text,
			Processed: text,
			Data:      resultMap,
		}, nil
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
		"sentiment":  sentiment,
		"score":      score,
		"confidence": confidence,
		"keywords":   keywords,
	}

	// Add debug info back if it existed
	if debugInfo != nil {
		resultMap["debug"] = debugInfo
	}

	return &Result{
		Original:  text,
		Processed: text,
		Data:      resultMap,
	}, nil
}

// ProcessItem processes a data.TextItem
func (p *SentimentProcessor) ProcessItem(ctx context.Context, item *data.TextItem) (*Result, error) {
	result, err := p.Process(ctx, item.Content)
	if err != nil {
		return nil, err
	}

	// Replace the generated original with the actual item
	result.Original = item

	return result, nil
}

// ProcessBatch processes a batch of items
func (p *SentimentProcessor) ProcessBatch(ctx context.Context, items []*data.TextItem) ([]*Result, error) {
	results := make([]*Result, len(items))

	for i, item := range items {
		result, err := p.ProcessItem(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessSource processes all items from a source
func (p *SentimentProcessor) ProcessSource(ctx context.Context, source data.Source, batchSize, workers int) ([]*Result, error) {
	processor := data.NewParallelProcessor(source, batchSize, workers)
	defer processor.Close()

	// Convert data.TextItem processor to Result processor
	itemProcessor := func(ctx context.Context, item *data.TextItem) (*data.TextItem, error) {
		result, err := p.ProcessItem(ctx, item)
		if err != nil {
			return nil, err
		}

		// Pack the result into the metadata
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		if item.Metadata == nil {
			item.Metadata = make(map[string]interface{})
		}
		item.Metadata["result"] = string(resultJSON)

		return item, nil
	}

	// Process all items
	processedItems, err := processor.ProcessAll(ctx, itemProcessor)
	if err != nil {
		return nil, err
	}

	// Extract results from metadata
	results := make([]*Result, len(processedItems))
	for i, item := range processedItems {
		resultJSON, ok := item.Metadata["result"].(string)
		if !ok {
			return nil, fmt.Errorf("missing result in item metadata")
		}

		var result Result
		if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
			return nil, err
		}

		results[i] = &result
	}

	return results, nil
}

// Register the processor with the registry
func init() {
	Register("sentiment", func(provider llm.Provider, options Options) (Processor, error) {
		return NewSentimentProcessor(provider, options)
	})
}
