package builtin

import (
	"github.com/eisenzopf/agentic-text/pkg/processor"
)

// AttributeMatch represents a match between required and available attributes
type AttributeMatch struct {
	// RequiredField is the name of the required attribute field
	RequiredField string `json:"required_field"`
	// MatchedField is the name of the matching available attribute
	MatchedField string `json:"matched_field"`
	// Confidence is the similarity confidence score (0.0-1.0)
	Confidence float64 `json:"confidence"`
	// MatchRationale explains why these attributes are considered similar
	MatchRationale string `json:"match_rationale"`
	// MatchType categorizes the type of match (exact, semantic, partial, etc.)
	MatchType string `json:"match_type"`
}

// MissingAttribute represents an attribute that has no suitable match
type MissingAttribute struct {
	// FieldName is the name of the missing attribute
	FieldName string `json:"field_name"`
	// Title is the human-readable title of the missing attribute
	Title string `json:"title"`
	// Description explains what this attribute represents
	Description string `json:"description"`
	// Reason explains why no match was found
	Reason string `json:"reason"`
	// Suggestions provide potential alternatives or workarounds
	Suggestions []string `json:"suggestions,omitempty"`
}

// MatchSummary provides an overview of the matching results
type MatchSummary struct {
	// TotalRequired is the number of required attributes
	TotalRequired int `json:"total_required"`
	// TotalMatched is the number of successfully matched attributes
	TotalMatched int `json:"total_matched"`
	// TotalMissing is the number of attributes with no suitable match
	TotalMissing int `json:"total_missing"`
	// MatchRate is the percentage of required attributes that were matched
	MatchRate float64 `json:"match_rate"`
	// AverageConfidence is the average confidence score of all matches
	AverageConfidence float64 `json:"average_confidence"`
	// Quality assessment of the overall matching results
	Quality string `json:"quality"`
}

// AttributeMatchResult contains comprehensive attribute matching results
type AttributeMatchResult struct {
	// Matches contains successful attribute matches
	Matches []AttributeMatch `json:"matches"`
	// MissingAttributes contains attributes that could not be matched
	MissingAttributes []MissingAttribute `json:"missing_attributes"`
	// MatchSummary provides statistical overview of matching results
	MatchSummary MatchSummary `json:"match_summary"`
	// Recommendations suggest next steps based on matching results
	Recommendations []string `json:"recommendations,omitempty"`
	// ProcessorType is the type of processor that generated this result
	ProcessorType string `json:"processor_type"`
}

// Register the processor with the registry
func init() {
	processor.NewBuilder("attribute_matcher").
		WithStruct(&AttributeMatchResult{}).
		WithContentTypes("text", "json").
		WithRole("You are an expert at semantic similarity analysis and attribute matching with deep understanding of data relationships and contextual meaning").
		WithObjective("Match required attributes against available attributes using semantic similarity, identify gaps, and provide detailed analysis of attribute relationships").
		WithInstructions(
			"Compare required attributes against available attributes to find semantic matches",
			"Consider field names, titles, descriptions, and conceptual meaning when matching",
			"Assign confidence scores based on semantic similarity and contextual relevance",
			"Identify match types: exact, semantic, partial, or conceptual matches",
			"Provide clear rationale explaining why attributes are considered similar",
			"For unmatched attributes, explain why no suitable match was found",
			"Suggest alternatives or workarounds for missing attributes",
			"Calculate match rates and provide quality assessment of overall results",
		).
		WithCustomSection("Matching Criteria", `
Match Types and Criteria:
- Exact Match (0.95-1.0): Identical or nearly identical field names and meanings
- Strong Semantic Match (0.8-0.94): Same concept with different terminology
- Moderate Match (0.6-0.79): Related concepts that capture similar information
- Weak Match (0.4-0.59): Loosely related but potentially useful
- No Match (0.0-0.39): No meaningful relationship

Consider:
- Field name similarity and common abbreviations
- Conceptual meaning and purpose
- Data type and structure compatibility
- Business context and domain relevance
- Synonyms and alternative terminology`).
		WithCustomSection("Confidence Assessment", `
Confidence Scoring Guidelines:
- 1.0: Perfect match, identical meaning and purpose
- 0.9-0.99: Excellent match, same concept with minor variations
- 0.8-0.89: Good match, captures the same essential information
- 0.7-0.79: Acceptable match, similar purpose with some differences
- 0.6-0.69: Marginal match, related but may miss some aspects
- Below 0.6: Poor match, significant differences in meaning

Factors affecting confidence:
- Semantic similarity of names and descriptions
- Conceptual alignment and purpose
- Data type and format compatibility
- Domain-specific meaning and context`).
		WithCustomSection("Gap Analysis Guidelines", `
For missing attributes:
- Clearly explain why no suitable match exists
- Identify the closest alternatives and their limitations
- Suggest potential workarounds or data transformations
- Recommend data collection or enhancement strategies
- Consider if multiple available attributes could combine to fulfill the requirement

Quality Assessment:
- Excellent (90%+ match rate, high confidence): Ready for implementation
- Good (70-89% match rate): Usable with minor gaps
- Fair (50-69% match rate): Significant gaps requiring attention
- Poor (<50% match rate): Major restructuring needed`).
		Register()
}
