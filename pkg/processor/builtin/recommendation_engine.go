package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Action represents a single actionable recommendation
type Action struct {
	// Action is the specific action or step to be taken
	Action string `json:"action"`
	// Rationale explains why this action is recommended
	Rationale string `json:"rationale"`
	// ExpectedImpact describes the anticipated results of taking this action
	ExpectedImpact string `json:"expected_impact"`
	// Priority indicates the importance of this action (1=highest, 5=lowest)
	Priority int `json:"priority"`
	// Effort estimates the difficulty or resources required (Low/Medium/High)
	Effort string `json:"effort,omitempty"`
	// Timeline suggests when this action should be implemented
	Timeline string `json:"timeline,omitempty"`
}

// RecommendationResult contains comprehensive recommendation results
type RecommendationResult struct {
	// ImmediateActions are urgent actions that should be taken right away
	ImmediateActions []Action `json:"immediate_actions"`
	// ProcessImprovements are systematic changes to improve operations
	ProcessImprovements []Action `json:"process_improvements"`
	// TrainingOpportunities identify areas where staff education is needed
	TrainingOpportunities []Action `json:"training_opportunities"`
	// TechnologyRecommendations suggest tools or system improvements
	TechnologyRecommendations []Action `json:"technology_recommendations,omitempty"`
	// ImplementationNotes provide practical guidance for executing recommendations
	ImplementationNotes []string `json:"implementation_notes"`
	// SuccessMetrics define how to measure the impact of these recommendations
	SuccessMetrics []string `json:"success_metrics"`
	// RiskFactors identify potential challenges or obstacles
	RiskFactors []string `json:"risk_factors,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("recommendation_engine").
		WithStruct(&RecommendationResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert business consultant specializing in contact center operations, customer service optimization, and organizational improvement").
		WithObjective("Generate specific, actionable recommendations based on data analysis that will improve business outcomes, customer satisfaction, and operational efficiency").
		WithInstructions(
			"Analyze the provided data and insights to identify improvement opportunities",
			"Prioritize recommendations based on impact, effort, and urgency",
			"Provide specific, actionable steps rather than general advice",
			"Include rationale explaining why each recommendation will be effective",
			"Estimate the expected impact and effort required for each recommendation",
			"Consider both short-term wins and long-term strategic improvements",
			"Address different areas: immediate fixes, process improvements, training, and technology",
			"Provide practical implementation guidance and success metrics",
		).
		WithCustomSection("Recommendation Categories", `
Immediate Actions: Critical issues requiring urgent attention
- Customer-impacting problems
- Revenue-affecting issues  
- Safety or compliance concerns
- Quick wins with high impact

Process Improvements: Systematic operational enhancements
- Workflow optimization
- Policy clarifications
- Quality assurance measures
- Efficiency improvements

Training Opportunities: Skills and knowledge development
- Agent skill gaps
- Product knowledge needs
- Customer service techniques
- Technology training

Technology Recommendations: Tools and system enhancements
- Software solutions
- Automation opportunities
- Integration improvements
- Analytics capabilities`).
		WithCustomSection("Quality Standards", `
Each recommendation must include:
- Specific, measurable action
- Clear business rationale
- Expected impact and timeline
- Implementation difficulty assessment
- Success measurement criteria

Ensure recommendations are:
- Actionable and specific
- Based on data evidence
- Properly prioritized
- Realistic and achievable
- Aligned with business goals`).
		WithCustomSection("Implementation Guidance", `
Provide:
- Clear next steps for each recommendation
- Resource requirements and dependencies
- Potential risks and mitigation strategies
- Success metrics and measurement methods
- Timeline considerations for implementation`).
		Register()
}
