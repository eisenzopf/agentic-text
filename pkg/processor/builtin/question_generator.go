package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// ResearchQuestion represents a single research question
type ResearchQuestion struct {
	// QuestionID is a unique identifier for this question
	QuestionID string `json:"question_id"`
	// Question is the actual research question text
	Question string `json:"question"`
	// Rationale explains why this question is important to ask
	Rationale string `json:"rationale"`
	// Priority indicates the importance of this question (1=highest, 5=lowest)
	Priority int `json:"priority"`
	// Category classifies the type of question (e.g., "operational", "strategic", "customer")
	Category string `json:"category"`
	// RequiredData identifies what data would be needed to answer this question
	RequiredData []string `json:"required_data,omitempty"`
	// ExpectedInsight describes what insights this question could provide
	ExpectedInsight string `json:"expected_insight,omitempty"`
}

// QuestionGenerationResult contains comprehensive question generation results
type QuestionGenerationResult struct {
	// Questions contains the generated research questions
	Questions []ResearchQuestion `json:"questions"`
	// Context provides background information that influenced question generation
	Context string `json:"context"`
	// TotalQuestions is the number of questions generated
	TotalQuestions int `json:"total_questions"`
	// Categories lists the different question categories represented
	Categories []string `json:"categories,omitempty"`
	// ResearchAreas identifies the main areas of inquiry
	ResearchAreas []string `json:"research_areas,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("question_generator").
		WithStruct(&QuestionGenerationResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert research methodologist and data analyst specializing in customer service and contact center research").
		WithObjective("Generate insightful, actionable research questions that will uncover valuable insights from conversation data and drive data-driven decision making").
		WithInstructions(
			"Analyze the provided context to understand the business domain and objectives",
			"Generate relevant research questions that would provide valuable insights",
			"Prioritize questions based on business impact and feasibility",
			"Categorize questions by research area (operational, strategic, customer experience, etc.)",
			"Provide clear rationale explaining why each question is important",
			"Identify what data would be required to answer each question",
			"Focus on questions that can lead to actionable insights and improvements",
			"Consider both immediate operational questions and strategic long-term inquiries",
		).
		WithCustomSection("Question Categories", `
Research Question Categories:
- Operational: Day-to-day performance and efficiency
- Strategic: Long-term planning and direction
- Customer Experience: Customer satisfaction and journey
- Quality Assurance: Service quality and standards
- Training: Skills development and knowledge gaps
- Process Improvement: Workflow and procedure optimization
- Technology: Tools and system effectiveness
- Business Impact: Revenue, cost, and ROI considerations
- Trend Analysis: Patterns and changes over time
- Competitive: Market position and differentiation`).
		WithCustomSection("Question Quality Criteria", `
Good Research Questions Should Be:
- Specific: Clear and well-defined scope
- Measurable: Can be answered with available or obtainable data
- Actionable: Results can inform decisions and improvements
- Relevant: Important to business objectives and stakeholders
- Time-bound: Consider temporal aspects and urgency

Avoid:
- Overly broad or vague questions
- Questions that can't be answered with available data
- Leading questions that assume conclusions
- Questions with obvious or trivial answers
- Multiple questions bundled into one`).
		WithCustomSection("Prioritization Guidelines", `
Priority Levels:
1. Critical: Urgent business needs, high impact decisions
2. High: Important strategic questions, significant improvement opportunities
3. Medium: Valuable insights, moderate business impact
4. Low: Interesting but not immediately critical
5. Future: Long-term research considerations

Consider:
- Business impact and urgency
- Data availability and analysis feasibility
- Stakeholder interest and needs
- Resource requirements for investigation
- Potential for actionable outcomes`).
		Register()
}
