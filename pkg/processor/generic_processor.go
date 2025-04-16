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
	validateStructure bool,
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
			validateStructure: validateStructure,
		}

		// Set the default responder
		responseHandler.DefaultResponder = func() interface{} {
			return responseHandler.createDefaultResponse()
		}

		// Configure fields based on the result struct
		responseHandler.configureFieldsFromStruct()

		// Apply processor-specific defaults
		responseHandler.applyProcessorDefaults()

		// Check for custom field validators (ValidateFieldName methods)
		// These run *after* the main structure validation (if enabled and passed)
		// Iterate over fields defined in the handler (which come from ResultStruct)
		for fieldName := range responseHandler.Fields {
			// Build the expected custom validator method name: "Validate" + Title case field name
			methodName := "Validate" + strings.Title(fieldName)
			validatorMethod := reflect.ValueOf(resultStruct).MethodByName(methodName)

			// If a custom validator method exists, add it to DynamicValidators
			if validatorMethod.IsValid() {
				// Call the validator method to get transform function
				results := validatorMethod.Call(nil)
				if len(results) > 0 {
					if transformFn, ok := results[0].Interface().(func(interface{}) interface{}); ok {
						// Add custom validator/transformer to DynamicValidators
						responseHandler.DynamicValidators[fieldName] = transformFn
					}
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
