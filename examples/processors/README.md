# Processor Testing Application

A comprehensive testing application for all available processors in the agentic-text framework, designed for easy evaluation and experimentation with conversation analysis capabilities.

## Overview

This application provides multiple ways to test and evaluate the 12 available processors using realistic customer service conversation data. It's inspired by the Python contact center analysis examples but built specifically for the Go framework.

## Features

### 🔬 **Comprehensive Testing**
- Test all 12 processors with realistic conversation data
- Support for batch testing across multiple processors
- Individual processor testing with specific conversations
- Custom text input for any processor

### 🎮 **Interactive Mode**
- Command-line interface for real-time testing
- Easy processor selection and conversation browsing
- Custom text input with multi-line support
- Built-in help and documentation

### 📊 **Rich Output**
- JSON-formatted results with timing information
- Verbose mode for debugging and inspection
- Success/failure tracking with detailed error reporting
- Export results to JSON files for analysis

### 💬 **Sample Conversations**
- 5 realistic customer service conversations
- Multiple categories: fee disputes, technical support, billing, cancellations, product inquiries
- Conversation summaries and categorization
- Real-world scenarios for comprehensive testing

## Available Processors

### Core Text Analysis (6 processors)
- `sentiment` - Sentiment analysis with scores and keywords
- `intent` - Intent classification for customer service
- `keyword_extraction` - Keyword extraction with relevance scoring
- `speech_act` - Speech act analysis (questions, requests, etc.)
- `get_attributes` - Structured attribute extraction
- `required_attributes` - Identify required attributes for analysis

### Advanced Analysis (6 processors)
- `data_analyzer` - Comprehensive data analysis and pattern identification
- `categorizer` - Advanced categorization with label consolidation
- `recommendation_engine` - Actionable recommendation generation
- `quality_reviewer` - LLM output quality assessment
- `attribute_matcher` - Semantic attribute matching
- `question_generator` - Research question generation

## Setup

1. **Set up API keys** (choose one):
   ```bash
   export OPENAI_API_KEY="your-openai-key"
   # OR
   export GEMINI_API_KEY="your-gemini-key"
   ```

2. **Navigate to the processors directory**:
   ```bash
   cd examples/processors
   ```

## Usage Examples

### 🎯 **Test Specific Processor**
Test a single processor with all sample conversations:
```bash
go run main.go -processor sentiment
go run main.go -processor data_analyzer
go run main.go -processor recommendation_engine
```

With verbose output for debugging:
```bash
go run main.go -processor sentiment -verbose
```

### 🧪 **Test All Processors**
Run comprehensive testing across all processors:
```bash
go run main.go -test-all
```

Save results to file:
```bash
go run main.go -test-all -output results.json
```

### 🎮 **Interactive Mode**
Launch interactive mode for exploration:
```bash
go run main.go -interactive
```

Interactive commands:
```
processor-test> help                    # Show available commands
processor-test> list                    # List all processors
processor-test> conversations           # Show sample conversations
processor-test> test sentiment          # Test sentiment with all conversations
processor-test> test intent conv_001    # Test intent with specific conversation
processor-test> custom data_analyzer    # Test with custom text input
processor-test> quit                    # Exit
```

### ✏️ **Custom Text Input**
Test any processor with your own text:
```bash
go run main.go -interactive
# Then use: custom <processor_name>
# Enter your text and end with "###"
```

## Sample Conversations

The application includes 5 realistic customer service conversations:

1. **Fee Dispute** (`conv_001`) - Customer disputing late payment fee
2. **Technical Support** (`conv_002`) - Internet connectivity troubleshooting  
3. **Billing Inquiry** (`conv_003`) - Customer confused about bill increase
4. **Service Cancellation** (`conv_004`) - Customer relocating and cancelling service
5. **Product Inquiry** (`conv_005`) - Customer interested in upgrading internet speed

Each conversation includes:
- **Full dialogue** between customer and agent
- **Category classification** for organization
- **Summary description** for quick reference
- **Realistic scenarios** reflecting actual customer service interactions

## Processor-Specific Input Preparation

The application intelligently prepares input for different processors:

- **`required_attributes`** - Provides research questions about the conversation
- **`data_analyzer`** - Structures conversation data with analysis questions
- **`recommendation_engine`** - Includes conversation context and category
- **`quality_reviewer`** - Formats content for quality assessment
- **`attribute_matcher`** - Provides required and available attribute examples
- **`question_generator`** - Includes business context for question generation
- **Others** - Uses conversation text directly

## Output Format

All results are returned in structured JSON format:

```json
{
  "processor_name": "sentiment",
  "input": "conv_001",
  "duration": "1.234s",
  "output": {
    "sentiment": "positive", 
    "score": 0.8,
    "confidence": 0.95,
    "keywords": ["thank you", "appreciate", "great customer"],
    "processor_type": "sentiment"
  }
}
```

Error handling:
```json
{
  "processor_name": "sentiment",
  "input": "conv_001", 
  "duration": "0.123s",
  "error": "API rate limit exceeded"
}
```

## Use Cases

### 🔍 **Research & Development**
- Test new processor implementations
- Compare processor performance across different conversation types
- Validate processor accuracy with known conversation categories

### 📈 **Performance Analysis**
- Measure processing times across different processors
- Identify performance bottlenecks in complex workflows
- Test scalability with batch processing

### 🛠️ **Debugging & Troubleshooting**
- Use verbose mode to inspect LLM prompts and responses
- Test edge cases with custom conversation inputs
- Validate processor behavior with specific conversation types

### 📚 **Learning & Education**
- Understand how different processors work with real conversation data
- Explore the capabilities of the agentic-text framework
- Learn best practices for conversation analysis

## Advanced Features

### **Batch Testing**
Test multiple processors efficiently:
```bash
go run main.go -test-all -verbose > detailed_results.log
```

### **Result Analysis**
Export and analyze results:
```bash
go run main.go -test-all -output results.json
# Then analyze results.json with your preferred tools
```

### **Custom Workflows**
Use interactive mode to build custom analysis workflows:
```
1. Generate research questions with question_generator
2. Identify required attributes with required_attributes  
3. Extract attributes with get_attributes
4. Analyze data with data_analyzer
5. Generate recommendations with recommendation_engine
6. Review quality with quality_reviewer
```

## Comparison with Python Examples

This Go implementation provides similar functionality to the Python contact center analysis examples:

| Python Feature | Go Implementation |
|----------------|------------------|
| `analyze_fee_disputes.py` | Comprehensive testing with `data_analyzer` and `recommendation_engine` |
| `generate_attributes.py` | `get_attributes` and `required_attributes` processors |
| `generate_intents.py` | `intent` and `categorizer` processors |
| Interactive analysis | Interactive mode with command-line interface |
| Batch processing | `-test-all` flag with multiple conversations |
| Error handling | Comprehensive error reporting and graceful failures |

## Contributing

To add new sample conversations:
1. Add to the `sampleConversations` slice in `main.go`
2. Include realistic customer-agent dialogue
3. Provide category and summary metadata
4. Test with multiple processors to ensure compatibility

To extend functionality:
1. Add new command-line flags for additional options
2. Implement new interactive commands
3. Add processor-specific input preparation logic
4. Enhance output formatting and analysis features 