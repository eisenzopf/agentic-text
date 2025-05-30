package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// AnalysisAnswer represents a single answer to a research question
type AnalysisAnswer struct {
	// Question is the research question being answered
	Question string `json:"question"`
	// Answer is the detailed response to the question
	Answer string `json:"answer"`
	// KeyMetrics are quantifiable metrics that support the answer
	KeyMetrics []string `json:"key_metrics"`
	// Confidence indicates the reliability of this answer (High/Medium/Low)
	Confidence string `json:"confidence"`
	// SupportingData provides evidence from the data that supports this answer
	SupportingData string `json:"supporting_data"`
}

// Pattern represents an identified pattern in the data
type Pattern struct {
	// Name is a short name for the pattern
	Name string `json:"name"`
	// Description explains what the pattern represents
	Description string `json:"description"`
	// Frequency indicates how often this pattern occurs
	Frequency string `json:"frequency"`
	// Significance explains why this pattern is important
	Significance string `json:"significance"`
}

// DataAnalysisResult contains comprehensive analysis results
type DataAnalysisResult struct {
	// Answers contains responses to the research questions
	Answers []AnalysisAnswer `json:"answers"`
	// DataGaps identifies limitations in the available data
	DataGaps []string `json:"data_gaps"`
	// KeyMetrics are the most important quantifiable findings
	KeyMetrics []string `json:"key_metrics"`
	// Patterns contains identified patterns in the data
	Patterns []Pattern `json:"patterns,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("data_analyzer").
		WithStruct(&DataAnalysisResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert data analyst specializing in contact center analytics and customer service research").
		WithObjective("Analyze customer service data to answer research questions, identify patterns, and provide actionable insights with supporting evidence").
		WithInstructions(
			"Analyze the provided data against the research questions",
			"Provide specific, detailed answers citing the data as evidence",
			"Identify key quantifiable metrics that support each answer",
			"Assess confidence levels (High/Medium/Low) based on data quality and sample size",
			"Identify any data gaps or limitations that affect the analysis",
			"Look for patterns and trends in the data that provide additional insights",
			"Ensure all answers are supported by concrete evidence from the dataset",
		).
		WithCustomSection("Analysis Guidelines", `
Focus on:
- Quantifiable insights with supporting evidence from the data
- Pattern identification across conversations and interactions
- Confidence assessment based on data quality, sample size, and consistency
- Clear identification of limitations, gaps, and areas needing more data
- Actionable insights that can drive business decisions
- Statistical significance and trend analysis where applicable`).
		WithCustomSection("Output Quality Standards", `
For each answer:
- Cite specific data points and statistics
- Explain the methodology used to reach conclusions
- Provide context about data limitations
- Include confidence levels with justification
- Suggest areas where additional data would improve accuracy`).
		Register()
}
