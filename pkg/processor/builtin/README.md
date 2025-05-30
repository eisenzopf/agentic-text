# Built-in Processors

This package provides a collection of ready-to-use processors for common text analysis tasks using the processor framework.

These processors are registered automatically when the package is imported, so you can immediately use them with `processor.Create()`

## Available Processors

### Core Text Analysis (Phase 0 - Existing)

#### `sentiment` - Sentiment Analysis
Analyzes the sentiment of text, returning sentiment type, score, confidence, and keywords.

**Output:** Sentiment classification with numerical scores and supporting keywords
**Use Cases:** Customer feedback analysis, social media monitoring, content evaluation

#### `intent` - Intent Classification  
Identifies the primary intent in customer service conversations.

**Output:** Structured intent classification with labels and descriptions
**Use Cases:** Customer service routing, chatbot training, conversation analysis

#### `keyword_extraction` - Keyword Extraction
Extracts important keywords from text with relevance scores and categories.

**Output:** List of keywords with relevance scores and semantic categories
**Use Cases:** Content tagging, search optimization, topic identification

#### `speech_act` - Speech Act Analysis
Identifies speech acts (questions, requests, statements, etc.) within text.

**Output:** Classification of speech acts with complexity scores and keywords
**Use Cases:** Conversation analysis, dialogue systems, communication research

#### `get_attributes` - Attribute Extraction
Extracts structured attributes and their values from text based on provided schemas.

**Output:** Structured attribute-value pairs with confidence scores
**Use Cases:** Information extraction, form filling, data structuring

#### `required_attributes` - Required Attribute Identification
Identifies data attributes required to answer specific questions.

**Output:** List of required attributes with descriptions and rationale
**Use Cases:** Research planning, data requirements analysis, schema design

### Advanced Analysis (Phase 1 - Core Analysis)

#### `data_analyzer` - Advanced Data Analysis
Analyzes customer service data to answer research questions and identify patterns.

**Output:** Comprehensive analysis with answers, patterns, and data gaps
**Use Cases:** Research question answering, pattern identification, business intelligence

#### `categorizer` - Advanced Categorization
Categorizes and classifies conversation elements with label consolidation.

**Output:** Classifications, label mappings, and semantic groupings
**Use Cases:** Content organization, label consolidation, semantic clustering

#### `recommendation_engine` - Action Recommendations
Generates actionable recommendations based on analysis results.

**Output:** Prioritized recommendations across multiple categories
**Use Cases:** Process improvement, strategic planning, operational optimization

### Quality & Matching (Phase 2 - Enhancement)

#### `quality_reviewer` - LLM Output Review
Reviews and evaluates LLM-generated content for quality and accuracy.

**Output:** Quality scores, improvement suggestions, and prompt effectiveness analysis
**Use Cases:** Quality assurance, content validation, prompt optimization

#### `attribute_matcher` - Semantic Matching
Matches required attributes against available attributes using semantic similarity.

**Output:** Attribute matches with confidence scores and gap analysis
**Use Cases:** Schema mapping, data integration, requirement analysis

#### `question_generator` - Research Question Generation
Generates research questions about conversation data and business contexts.

**Output:** Prioritized research questions with rationale and data requirements
**Use Cases:** Research planning, business analysis, insight discovery

## Usage Examples

### Basic Usage

```go
import (
    "github.com/eisenzopf/agentic-text/pkg/processor"
    _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin" // Import for registration
)

// Create any processor
sentimentProc, err := processor.Create("sentiment", provider, options)
analyzerProc, err := processor.Create("data_analyzer", provider, options)
```

### Using with Easy Library

```go
import (
    "github.com/eisenzopf/agentic-text/pkg/easy"
    _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
)

// All processors are available through easy library
result, err := easy.ProcessText("customer feedback text", "sentiment")
analysis, err := easy.ProcessText("research data", "data_analyzer")
recommendations, err := easy.ProcessText("analysis results", "recommendation_engine")
```

### Advanced Workflows

```go
// Multi-step analysis pipeline
questions := "What are the main customer pain points?"

// 1. Generate required attributes
attributes, err := easy.ProcessText(questions, "required_attributes")

// 2. Extract attributes from conversations  
extracted, err := easy.ProcessText(conversations, "get_attributes")

// 3. Analyze the data
analysis, err := easy.ProcessText(extracted, "data_analyzer")

// 4. Generate recommendations
recommendations, err := easy.ProcessText(analysis, "recommendation_engine")

// 5. Review quality
quality, err := easy.ProcessText(analysis, "quality_reviewer")
```

## Processor Categories

### By Use Case
- **Customer Service:** sentiment, intent, speech_act, recommendation_engine
- **Research & Analysis:** data_analyzer, question_generator, required_attributes
- **Data Processing:** get_attributes, attribute_matcher, categorizer
- **Quality Assurance:** quality_reviewer, keyword_extraction

### By Input Type
- **Text Analysis:** sentiment, intent, speech_act, keyword_extraction
- **Structured Data:** data_analyzer, categorizer, attribute_matcher
- **Meta-Analysis:** quality_reviewer, question_generator, recommendation_engine

### By Output Complexity
- **Simple:** sentiment, intent, keyword_extraction
- **Structured:** get_attributes, required_attributes, speech_act
- **Comprehensive:** data_analyzer, recommendation_engine, quality_reviewer

## Implementation Notes

All processors are implemented using the new builder pattern for:
- **Minimal code** (67% reduction vs manual implementation)
- **Consistent structure** across all processors
- **Easy maintenance** and updates
- **Flexible customization** when needed

See individual processor files for detailed implementation and customization options.

## Content Type Support

Most processors support both:
- `text` - Plain text input
- `json` - Structured JSON input

Check individual processor documentation for specific content type requirements.

## Importing

To use these processors, simply import this package for its side effects:

```go
import _ "github.com/eisenzopf/agentic-text/pkg/processor/builtin"
```

This will register all builtin processors with the processor registry, making them available through `processor.Create()`. 