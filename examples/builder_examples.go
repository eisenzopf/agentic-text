package examples

import (
	"context"
	"fmt"

	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Example 1: MINIMAL - Just struct + name (auto-generates everything)
type SimpleSentimentResult struct {
	Sentiment     string  `json:"sentiment"`
	Score         float64 `json:"score"`
	ProcessorType string  `json:"processor_type"`
}

func registerMinimalSentiment() {
	processor.NewBuilder("minimal_sentiment").
		WithStruct(&SimpleSentimentResult{}).
		Register()
}

// Example 2: BASIC - Add role and objective
type BasicSentimentResult struct {
	Sentiment     string   `json:"sentiment"`
	Score         float64  `json:"score"`
	Confidence    float64  `json:"confidence"`
	Keywords      []string `json:"keywords"`
	ProcessorType string   `json:"processor_type"`
}

func registerBasicSentiment() {
	processor.NewBuilder("basic_sentiment").
		WithStruct(&BasicSentimentResult{}).
		WithRole("You are an expert sentiment analysis tool that ONLY outputs valid JSON").
		WithObjective("Analyze sentiment considering tone, word choices, context, and nuances like sarcasm").
		Register()
}

// Example 3: DETAILED - Add specific instructions
func registerDetailedSentiment() {
	processor.NewBuilder("detailed_sentiment").
		WithStruct(&BasicSentimentResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert sentiment analysis tool that ONLY outputs valid JSON").
		WithObjective("Analyze sentiment expressed in text accurately and objectively. Consider overall tone, specific word choices, context, and potential nuances like sarcasm or mixed feelings").
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

// Example 4: CUSTOM SECTIONS - Add domain-specific guidance
type CustomerServiceSentimentResult struct {
	Sentiment        string   `json:"sentiment"`
	Score            float64  `json:"score"`
	Confidence       float64  `json:"confidence"`
	Keywords         []string `json:"keywords"`
	EscalationNeeded bool     `json:"escalation_needed"`
	ProcessorType    string   `json:"processor_type"`
}

func registerCustomerServiceSentiment() {
	processor.NewBuilder("customer_service_sentiment").
		WithStruct(&CustomerServiceSentimentResult{}).
		WithRole("You are a customer service sentiment analysis expert").
		WithObjective("Analyze customer sentiment in service interactions and determine if escalation is needed").
		WithInstructions(
			"Analyze the customer's emotional state from their message",
			"Consider urgency indicators, frustration levels, and satisfaction signals",
			"Determine if the situation requires escalation to a human agent",
		).
		WithCustomSection("Escalation Criteria", `
Set escalation_needed to true if any of these conditions are met:
- Extremely negative sentiment (score < -0.7)
- Threats or legal language detected
- Multiple unresolved issues mentioned
- Explicit request for manager/supervisor
- Indicates they will cancel service/leave`).
		WithCustomSection("Customer Service Context", `
This analysis is for a customer service platform. Focus on:
- Customer satisfaction indicators
- Pain points and frustrations  
- Urgency and priority levels
- Professional communication needs`).
		Register()
}

// Example 5: FULLY CUSTOM - Use custom prompt generator
type AdvancedSentimentResult struct {
	Sentiment       string             `json:"sentiment"`
	Score           float64            `json:"score"`
	Confidence      float64            `json:"confidence"`
	Emotions        map[string]float64 `json:"emotions"`
	SarcasmDetected bool               `json:"sarcasm_detected"`
	ProcessorType   string             `json:"processor_type"`
}

type CustomSentimentPrompt struct{}

func (p *CustomSentimentPrompt) GeneratePrompt(ctx context.Context, text string) (string, error) {
	return fmt.Sprintf(`
ADVANCED SENTIMENT ANALYSIS SYSTEM v2.0

MISSION: Perform sophisticated multi-dimensional sentiment analysis

TARGET TEXT: %s

ANALYSIS REQUIREMENTS:
• Primary sentiment classification
• Numerical scoring with high precision
• Multi-emotion detection (joy, anger, fear, sadness, surprise, disgust)
• Sarcasm and irony detection
• Confidence assessment

OUTPUT: JSON ONLY - NO OTHER TEXT
`, text), nil
}

func registerAdvancedSentiment() {
	processor.NewBuilder("advanced_sentiment").
		WithStruct(&AdvancedSentimentResult{}).
		WithCustomPrompt(&CustomSentimentPrompt{}).
		WithValidation().
		Register()
}

// Example 6: WITH CUSTOM INIT - Add special initialization logic
func registerSentimentWithInit() {
	processor.NewBuilder("sentiment_with_init").
		WithStruct(&BasicSentimentResult{}).
		WithRole("Expert sentiment analyzer").
		WithObjective("Analyze sentiment with high accuracy").
		WithCustomInit(func(p *processor.GenericProcessor) error {
			// Custom initialization logic
			fmt.Println("Initializing custom sentiment processor with special configuration")
			// Could set up caches, load models, configure clients, etc.
			return nil
		}).
		Register()
}

// Example 7: EXPERIMENTAL - Using all features
func registerExperimentalSentiment() {
	processor.NewBuilder("experimental_sentiment").
		WithStruct(&AdvancedSentimentResult{}).
		WithContentTypes("text", "json", "html").
		WithRole("You are an experimental AI system with advanced sentiment analysis capabilities").
		WithObjective("Perform cutting-edge sentiment analysis using latest psychological and linguistic research").
		WithInstructions(
			"Apply multi-layered sentiment analysis techniques",
			"Consider cultural and contextual factors",
			"Use advanced emotion detection algorithms",
			"Detect subtle linguistic patterns like sarcasm, irony, and implicit sentiment",
		).
		WithCustomSection("Experimental Features", `
This is an experimental processor. Apply these advanced techniques:
- Contextual embedding analysis
- Temporal sentiment shifts within text
- Implicit vs explicit sentiment detection
- Cross-cultural sentiment interpretation`).
		WithCustomSection("Quality Assurance", `
Ensure high-quality analysis by:
- Double-checking edge cases
- Validating emotional complexity
- Confirming cultural sensitivity`).
		WithValidation().
		WithCustomInit(func(p *processor.GenericProcessor) error {
			fmt.Println("Loading experimental sentiment models...")
			return nil
		}).
		Register()
}
