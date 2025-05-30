# Example Testing Session

This document shows example outputs and usage patterns for the processor testing application.

## Basic Processor Testing

### Testing Sentiment Analysis
```bash
$ go run main.go -processor sentiment

🔬 Processor Testing Application
=================================
Available processors: sentiment, intent, keyword_extraction, speech_act, get_attributes, required_attributes, data_analyzer, categorizer, recommendation_engine, quality_reviewer, attribute_matcher, question_generator
Sample conversations: 5

🎯 Testing processor: sentiment
--------------------------------------------------
Conversation 1 (fee_dispute): ✅ Success (1.234s)
Conversation 2 (technical_support): ✅ Success (1.156s)
Conversation 3 (billing_inquiry): ✅ Success (1.089s)
Conversation 4 (service_cancellation): ✅ Success (1.201s)
Conversation 5 (product_inquiry): ✅ Success (1.178s)
```

### Testing with Verbose Output
```bash
$ go run main.go -processor sentiment -verbose

🎯 Testing processor: sentiment
--------------------------------------------------
Conversation 1 (fee_dispute): ✅ Success (1.234s)

    Result:
{
  "sentiment": "positive",
  "score": 0.82,
  "confidence": 0.95,
  "keywords": ["thank you", "appreciate", "great customer", "welcome"],
  "processor_type": "sentiment"
}

Conversation 2 (technical_support): ✅ Success (1.156s)

    Result:
{
  "sentiment": "positive", 
  "score": 0.78,
  "confidence": 0.91,
  "keywords": ["thank you", "perfect", "amazing", "help"],
  "processor_type": "sentiment"
}
```

## Interactive Mode Examples

### Starting Interactive Mode
```bash
$ go run main.go -interactive

🔬 Processor Testing Application
=================================
Available processors: sentiment, intent, keyword_extraction, speech_act, get_attributes, required_attributes, data_analyzer, categorizer, recommendation_engine, quality_reviewer, attribute_matcher, question_generator
Sample conversations: 5

🎮 Interactive Mode
Type 'help' for commands, 'quit' to exit

processor-test> help
Available commands:
  help                     - Show this help message
  list                     - List available processors
  conversations, convs     - List sample conversations
  test <processor> [conv]  - Test processor with conversation(s)
  custom <processor>       - Test processor with custom text
  quit, exit               - Exit interactive mode

Examples:
  test sentiment           - Test sentiment with all conversations
  test intent conv_001     - Test intent with specific conversation
  custom data_analyzer     - Test data_analyzer with custom input
```

### Listing Available Resources
```bash
processor-test> list
Available processors:
  1. sentiment
  2. intent
  3. keyword_extraction
  4. speech_act
  5. get_attributes
  6. required_attributes
  7. data_analyzer
  8. categorizer
  9. recommendation_engine
  10. quality_reviewer
  11. attribute_matcher
  12. question_generator

processor-test> conversations
Sample conversations:
  1. fee_dispute (conv_001)
     Customer disputing late payment fee
  2. technical_support (conv_002)
     Internet connectivity issues troubleshooting
  3. billing_inquiry (conv_003)
     Customer confused about sudden bill increase
  4. service_cancellation (conv_004)
     Customer wanting to cancel service due to relocation
  5. product_inquiry (conv_005)
     Customer interested in upgrading internet speed
```

### Testing Specific Combinations
```bash
processor-test> test intent conv_001

Testing intent with conversation conv_001...
✅ Success (1.423s)
    Result:
{
  "intent": "billing_dispute",
  "confidence": 0.93,
  "description": "Customer is disputing a charge on their bill",
  "category": "billing",
  "keywords": ["fee", "dispute", "charge", "bill"],
  "processor_type": "intent"
}

processor-test> test data_analyzer conv_002

Testing data_analyzer with conversation conv_002...
✅ Success (2.847s)
    Result:
{
  "answers": [
    {
      "question": "What was the primary customer issue in this conversation?",
      "answer": "The customer was experiencing slow internet speeds, getting only 10 Mbps instead of the expected 100 Mbps.",
      "key_metrics": ["10 Mbps actual speed", "100 Mbps expected speed", "5-year-old router"],
      "confidence": "High",
      "supporting_data": "Customer explicitly stated speeds and router age"
    }
  ],
  "data_gaps": [],
  "key_metrics": ["Internet speed discrepancy", "Router age factor", "Resolution success"],
  "patterns": [
    {
      "name": "Hardware Upgrade Resolution",
      "description": "Technical issues resolved through equipment upgrade",
      "frequency": "Common",
      "significance": "Proactive customer retention strategy"
    }
  ],
  "processor_type": "data_analyzer"
}
```

### Custom Text Input
```bash
processor-test> custom sentiment
Enter your text (end with '###' on a new line):
Customer called very upset about being charged twice for the same service. 
Agent apologized and immediately processed a full refund.
Customer was satisfied with the quick resolution.
###

Testing sentiment with custom text...
✅ Success (1.156s)
    Result:
{
  "sentiment": "neutral_to_positive", 
  "score": 0.12,
  "confidence": 0.87,
  "keywords": ["upset", "apologized", "refund", "satisfied", "quick resolution"],
  "processor_type": "sentiment"
}
```

## Comprehensive Testing

### Testing All Processors
```bash
$ go run main.go -test-all

🧪 Testing all processors with all sample conversations...

[1/12] Testing processor: sentiment
--------------------------------------------------
  Conversation 1 (fee_dispute): ✅ Success (1.234s)
  Conversation 2 (technical_support): ✅ Success (1.156s)
  Conversation 3 (billing_inquiry): ✅ Success (1.089s)
  Conversation 4 (service_cancellation): ✅ Success (1.201s)
  Conversation 5 (product_inquiry): ✅ Success (1.178s)

[2/12] Testing processor: intent
--------------------------------------------------
  Conversation 1 (fee_dispute): ✅ Success (1.423s)
  Conversation 2 (technical_support): ✅ Success (1.367s)
  Conversation 3 (billing_inquiry): ✅ Success (1.298s)
  Conversation 4 (service_cancellation): ✅ Success (1.445s)
  Conversation 5 (product_inquiry): ✅ Success (1.389s)

[3/12] Testing processor: data_analyzer
--------------------------------------------------
  Conversation 1 (fee_dispute): ✅ Success (2.847s)
  Conversation 2 (technical_support): ✅ Success (2.756s)
  Conversation 3 (billing_inquiry): ✅ Success (2.934s)
  Conversation 4 (service_cancellation): ✅ Success (2.821s)
  Conversation 5 (product_inquiry): ✅ Success (2.689s)

... (additional processors) ...

📊 Test Summary: 60/60 successful (100.0%)
```

### Exporting Results
```bash
$ go run main.go -test-all -output results.json
# Results saved to results.json

$ head -20 results.json
[
  {
    "processor_name": "sentiment",
    "input": "conv_001",
    "duration": "1.234s",
    "output": {
      "sentiment": "positive",
      "score": 0.82,
      "confidence": 0.95,
      "keywords": ["thank you", "appreciate", "great customer"],
      "processor_type": "sentiment"
    }
  },
  {
    "processor_name": "sentiment",
    "input": "conv_002",
    "duration": "1.156s",
...
```

## Advanced Usage Patterns

### Workflow Testing
Use interactive mode to simulate a complete analysis workflow:

```bash
processor-test> test question_generator
# Generate research questions about customer service

processor-test> test required_attributes 
# Identify what attributes are needed

processor-test> test get_attributes
# Extract those attributes from conversations

processor-test> test data_analyzer
# Analyze the extracted data

processor-test> test recommendation_engine  
# Generate actionable recommendations

processor-test> test quality_reviewer
# Review the quality of all outputs
```

### Performance Comparison
```bash
# Test processing speed across different processors
$ go run main.go -test-all | grep "Success" | sort -k3
  Conversation 1 (fee_dispute): ✅ Success (0.892s)  # keyword_extraction
  Conversation 1 (fee_dispute): ✅ Success (1.234s)  # sentiment  
  Conversation 1 (fee_dispute): ✅ Success (1.423s)  # intent
  Conversation 1 (fee_dispute): ✅ Success (2.847s)  # data_analyzer
  Conversation 1 (fee_dispute): ✅ Success (3.156s)  # recommendation_engine
```

### Error Handling Examples
```bash
processor-test> test invalid_processor
Unknown processor: invalid_processor

# API rate limiting example
processor-test> test sentiment
Testing sentiment with conversation conv_001...
❌ Error: API rate limit exceeded, please try again later

# Network timeout example  
processor-test> test data_analyzer
Testing data_analyzer with conversation conv_001...
❌ Error: context deadline exceeded
```

## Tips for Effective Testing

1. **Start with simple processors** like `sentiment` or `intent` to verify basic functionality
2. **Use verbose mode** when debugging or learning how processors work
3. **Test with custom text** to validate processors with your specific use cases
4. **Export results** for analysis and comparison across different runs
5. **Use interactive mode** for exploratory testing and workflow development
6. **Monitor processing times** to understand performance characteristics

This testing application provides a comprehensive way to evaluate and understand the capabilities of all available processors in the agentic-text framework. 