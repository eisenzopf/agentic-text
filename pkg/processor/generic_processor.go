package processor

import (
	"context"
	"reflect"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// GenericProcessor is a type for processors that use the standard response handling pattern
type GenericProcessor struct {
	// Embed BaseProcessor to inherit all methods
	BaseProcessor
	// ResultStruct is the struct that will be used for results
	ResultStruct interface{}
	// responseHandler is the configured response handler for this processor
	responseHandler ResponseHandler
}

// HandleResponse implements ResponseHandler interface - handles the LLM response
func (p *GenericProcessor) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	// The response handler is now set directly in RegisterGenericProcessor
	// This method is kept for backward compatibility
	if p.responseHandler != nil {
		return p.responseHandler.HandleResponse(ctx, text, responseData)
	}

	// Fallback to creating a new handler
	handler := NewResponseHandler(p.name, p.ResultStruct)
	return handler.AutoProcessResponse(ctx, text, responseData)
}

// RegisterGenericProcessor creates and registers a processor with standard behavior
func RegisterGenericProcessor(
	name string,
	contentTypes []string,
	resultStruct interface{},
	promptGenerator PromptGenerator,
	customInit func(*GenericProcessor) error,
	validationOptions ...map[string]interface{},
) {
	// Register the processor creator function
	Register(name, func(provider llm.Provider, options Options) (Processor, error) {
		// Create a new generic processor
		p := &GenericProcessor{
			ResultStruct: resultStruct,
		}

		// Create client from provider
		client := llm.NewProviderClient(provider)

		// Create response handler with dynamic validators if needed
		responseHandler := &BaseResponseHandler{
			ProcessorType:     name,
			ResultStruct:      resultStruct,
			Fields:            make(map[string]FieldMapper),
			DynamicValidators: make(map[string]func(interface{}) interface{}),
		}

		// Set the default responder
		responseHandler.DefaultResponder = func() interface{} {
			return responseHandler.createDefaultResponse()
		}

		// Configure fields based on the result struct
		responseHandler.configureFieldsFromStruct()

		// Apply processor-specific defaults
		responseHandler.applyProcessorDefaults()

		// Apply validation options if provided
		if len(validationOptions) > 0 && validationOptions[0] != nil {
			// Get the struct value
			structVal := reflect.ValueOf(resultStruct).Elem()

			// Check if the struct has Validate methods for fields that need validation
			if fieldName, ok := validationOptions[0]["field_name"].(string); ok {
				// Build the validator method name: "Validate" + Title case field name
				methodName := "Validate" + strings.Title(fieldName)
				validatorMethod := reflect.ValueOf(resultStruct).MethodByName(methodName)

				// If there's no existing validator method but we want validation
				if !validatorMethod.IsValid() {
					// Get default value from options or create a default
					var defaultValue interface{}
					if val, ok := validationOptions[0]["default_value"]; ok {
						defaultValue = val
					} else {
						// Create a default value based on field type
						// Look for the field in the struct
						for i := 0; i < structVal.Type().NumField(); i++ {
							field := structVal.Type().Field(i)

							// Get JSON tag name
							tag := field.Tag.Get("json")
							if tag == "" {
								tag = strings.ToLower(field.Name)
							} else {
								tag = strings.Split(tag, ",")[0]
							}

							if tag == fieldName {
								// Create a default value of the appropriate type
								defaultValue = reflect.New(field.Type).Elem().Interface()
								break
							}
						}
					}

					// Add validate method to dynamic validators
					responseHandler.DynamicValidators[fieldName] = ValidateData(fieldName, defaultValue)
				}
			}
		}

		// Override the generic HandleResponse method to use our configured handler
		p.responseHandler = responseHandler

		// Create and embed base processor with the appropriate content types
		base := NewBaseProcessor(name, contentTypes, client, nil, promptGenerator, p.responseHandler, options)
		p.BaseProcessor = *base

		// Call custom initializer if provided
		if customInit != nil {
			if err := customInit(p); err != nil {
				return nil, err
			}
		}

		return p, nil
	})
}
