package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// CriteriaScore represents an evaluation score for a specific criterion
type CriteriaScore struct {
	// Criterion is the specific quality criteria being evaluated
	Criterion string `json:"criterion"`
	// Score is the numerical score for this criterion (0.0-1.0)
	Score float64 `json:"score"`
	// Assessment provides detailed evaluation commentary
	Assessment string `json:"assessment"`
	// ImprovementNeeded indicates if this area requires enhancement
	ImprovementNeeded bool `json:"improvement_needed"`
	// Suggestions provide specific ways to improve this criterion
	Suggestions []string `json:"suggestions,omitempty"`
}

// QualityAssessment provides an overall quality evaluation
type QualityAssessment struct {
	// Score is the overall quality score (0.0-1.0)
	Score float64 `json:"score"`
	// Grade is a letter grade representation (A, B, C, D, F)
	Grade string `json:"grade"`
	// Strengths identifies what the output does well
	Strengths []string `json:"strengths"`
	// Weaknesses identifies areas needing improvement
	Weaknesses []string `json:"weaknesses"`
	// Summary provides an overall assessment
	Summary string `json:"summary"`
}

// PromptReview evaluates the effectiveness of the original prompt
type PromptReview struct {
	// Assessment evaluates how well the prompt worked
	Assessment string `json:"assessment"`
	// Clarity assesses how clear and unambiguous the prompt was
	Clarity float64 `json:"clarity"`
	// Completeness evaluates if the prompt provided sufficient guidance
	Completeness float64 `json:"completeness"`
	// SuggestedImprovements provides specific prompt enhancement ideas
	SuggestedImprovements []string `json:"suggested_improvements"`
}

// Improvement represents a specific improvement suggestion
type Improvement struct {
	// Issue describes the problem or opportunity
	Issue string `json:"issue"`
	// Suggestion provides a specific improvement recommendation
	Suggestion string `json:"suggestion"`
	// Priority indicates importance (1=highest, 5=lowest)
	Priority int `json:"priority"`
	// Category classifies the type of improvement needed
	Category string `json:"category"`
	// Impact describes the expected benefit of this improvement
	Impact string `json:"impact,omitempty"`
}

// ReviewResult contains comprehensive quality review results
type ReviewResult struct {
	// CriteriaScores contains detailed evaluation against specific criteria
	CriteriaScores []CriteriaScore `json:"criteria_scores"`
	// OverallQuality provides a comprehensive quality assessment
	OverallQuality QualityAssessment `json:"overall_quality"`
	// PromptEffectiveness evaluates the original prompt that generated the content
	PromptEffectiveness PromptReview `json:"prompt_effectiveness"`
	// Improvements contains prioritized suggestions for enhancement
	Improvements []Improvement `json:"improvements"`
	// RecommendedActions provide specific next steps
	RecommendedActions []string `json:"recommended_actions"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("quality_reviewer").
		WithStruct(&ReviewResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert quality assurance specialist and content reviewer with deep expertise in evaluating LLM-generated content for accuracy, completeness, and usefulness").
		WithObjective("Evaluate LLM-generated content against quality criteria, identify improvement opportunities, and provide specific recommendations for enhancement").
		WithInstructions(
			"Evaluate the provided LLM output against the specified quality criteria",
			"Provide numerical scores (0.0-1.0) for each evaluation criterion",
			"Identify specific strengths and weaknesses in the content",
			"Assess the effectiveness of the original prompt in generating quality output",
			"Provide prioritized, actionable suggestions for improvement",
			"Consider accuracy, completeness, clarity, usefulness, and relevance",
			"Evaluate whether the output meets its intended purpose",
			"Suggest specific improvements to both content and prompting approach",
		).
		WithCustomSection("Quality Evaluation Criteria", `
Standard Evaluation Criteria:
- Accuracy: Factual correctness and reliability
- Completeness: Coverage of all requested aspects
- Clarity: Clear, understandable communication
- Relevance: Appropriateness to the context and purpose
- Usefulness: Practical value and actionability
- Structure: Logical organization and formatting
- Specificity: Concrete details vs. vague generalities
- Evidence: Support for claims and conclusions

Custom criteria may be provided for specific use cases.`).
		WithCustomSection("Assessment Guidelines", `
Scoring Scale:
- 0.9-1.0: Excellent - Exceeds expectations
- 0.8-0.89: Good - Meets expectations well
- 0.7-0.79: Satisfactory - Meets basic expectations
- 0.6-0.69: Needs Improvement - Below expectations
- 0.0-0.59: Poor - Significant deficiencies

Grade Mapping:
- A: 0.9-1.0 (Excellent)
- B: 0.8-0.89 (Good)  
- C: 0.7-0.79 (Satisfactory)
- D: 0.6-0.69 (Needs Improvement)
- F: 0.0-0.59 (Poor)`).
		WithCustomSection("Improvement Prioritization", `
Priority Levels:
1. Critical: Issues that make content unusable or misleading
2. High: Significant gaps that impact effectiveness
3. Medium: Improvements that would enhance quality
4. Low: Minor enhancements and polish
5. Optional: Nice-to-have improvements

Improvement Categories:
- Content: Substance and information quality
- Structure: Organization and flow
- Clarity: Communication effectiveness
- Accuracy: Factual correctness
- Completeness: Coverage gaps
- Prompt: Original prompt improvements`).
		Register()
}
