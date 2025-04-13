# Agentic Text Processors

This package contains processors for analyzing text using LLMs.

## Generic Validation

The package now supports a generic validation approach that can be used to validate data returned from LLM responses without needing to implement custom validation methods for each processor.

### Using Generic Validation

When registering a new processor, you can specify validation options to automatically validate specific fields:

```go
// Register a processor with validation for the "attributes" field
RegisterGenericProcessor(
    "my_processor",         // name
    []string{"text", "json"}, // contentTypes
    &MyResultStruct{},       // resultStruct
    &MyPromptGenerator{},    // promptGenerator
    nil,                     // no custom initialization needed
    map[string]interface{}{  // validation options
        "field_name": "attributes",
        "default_value": []MyAttribute{
            {
                ID: "default",
                Value: "Default value used when validation fails",
            },
        },
    },
)
```

### Validation Options

The validation options map supports the following keys:

- `field_name`: The name of the field to validate (required)
- `default_value`: The default value to return if validation fails (optional)

### How It Works

The validation system will:

1. Check if the field exists in the LLM response
2. Validate that the field has the expected structure
3. If validation fails, return the default value
4. Apply the validation before any other transformations

### Custom Validation

You can still implement custom validation by adding a method to your result struct with the format:

```go
func (r *MyResultStruct) ValidateFieldName() func(interface{}) interface{} {
    // Custom validation logic
}
```

Where `FieldName` is the name of the field to validate (with first letter capitalized). 