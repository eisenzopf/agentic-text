package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// Classification represents a single classification result
type Classification struct {
	// Item is the text or element being classified
	Item string `json:"item"`
	// Category is the assigned category or class
	Category string `json:"category"`
	// IsMatch indicates if the item matches the target classification criteria
	IsMatch bool `json:"is_match"`
	// Confidence is the confidence score for this classification (0.0-1.0)
	Confidence float64 `json:"confidence"`
	// Rationale explains why this classification was assigned
	Rationale string `json:"rationale"`
}

// LabelGroup represents a group of semantically similar labels
type LabelGroup struct {
	// Theme is the common theme or concept that unites this group
	Theme string `json:"theme"`
	// Items are the individual labels or items in this group
	Items []string `json:"items"`
	// Rationale explains why these items were grouped together
	Rationale string `json:"rationale"`
	// Frequency indicates how common this group is in the data
	Frequency int `json:"frequency,omitempty"`
}

// CategorizationResult contains categorization and classification results
type CategorizationResult struct {
	// Classifications contains individual classification results
	Classifications []Classification `json:"classifications"`
	// LabelMapping maps original labels to consolidated/grouped labels
	LabelMapping map[string]string `json:"label_mapping,omitempty"`
	// Groups contains semantically similar items grouped together
	Groups []LabelGroup `json:"groups,omitempty"`
	// Summary provides an overview of the categorization results
	Summary string `json:"summary,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("categorizer").
		WithStruct(&CategorizationResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert at categorizing and classifying text content with advanced semantic understanding").
		WithObjective("Categorize items, classify content against criteria, consolidate similar labels, and create meaningful semantic groups").
		WithInstructions(
			"Analyze each item for classification against the specified criteria or categories",
			"Provide clear rationale for each classification decision",
			"Assign confidence scores based on how clearly the item fits the criteria",
			"When consolidating labels, group semantically similar items together",
			"Create meaningful themes that capture the essence of each group",
			"Maintain consistency in classification criteria across all items",
			"Consider context and domain-specific meanings when categorizing",
		).
		WithCustomSection("Classification Guidelines", `
Classification Criteria:
- Use semantic similarity and meaning, not just keyword matching
- Consider context and domain-specific interpretations
- Provide confidence scores reflecting classification certainty
- Explain rationale with specific reasons for each decision
- Group similar concepts even if expressed differently
- Maintain consistency across the entire dataset`).
		WithCustomSection("Label Consolidation Rules", `
When consolidating labels:
- Group synonyms and semantically equivalent terms
- Use the most clear and representative term as the group theme
- Preserve important distinctions while reducing redundancy
- Consider frequency and business importance when choosing representative terms
- Explain consolidation decisions with clear rationale`).
		WithCustomSection("Quality Standards", `
Ensure:
- Consistent classification criteria application
- Clear and actionable group themes
- Balanced consolidation (neither too granular nor too broad)
- Preservation of important semantic distinctions
- Business-relevant categorizations that enable action`).
		Register()
}
