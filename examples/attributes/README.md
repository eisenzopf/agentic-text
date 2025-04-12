# Attributes Example

This example demonstrates how to extract structured attributes from conversations using the agentic-text library. It works in two steps:

1. Use the `required_attributes` processor to determine what attributes we need to extract from a set of questions
2. Use the `get_attributes` processor to extract those attributes from actual conversation data

## Use Case

This example simulates analyzing customer cancellation calls to extract key information:
- Reasons for cancellation
- Whether agents attempt to save the customer
- What save offers are provided and which ones work

## How to Run

```
go run main.go
```

## How It Works

1. The example starts by defining questions about what information we want to extract
2. The `required_attributes` processor analyzes these questions and determines necessary data fields
3. A sample customer service conversation is provided
4. The `get_attributes` processor extracts the defined attributes from the conversation
5. The results are displayed in JSON format

This approach enables automated extraction of structured data from unstructured conversations. 