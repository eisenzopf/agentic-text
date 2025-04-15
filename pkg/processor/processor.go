package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/eisenzopf/agentic-text/pkg/data"
	"github.com/eisenzopf/agentic-text/pkg/llm"
)

// Options holds common configuration for processors
type Options struct {
	// PreProcessOptions holds options for pre-processing
	PreProcessOptions map[string]interface{}
	// LLMOptions holds options for LLM processing
	LLMOptions map[string]interface{}
	// PostProcessOptions holds options for post-processing
	PostProcessOptions map[string]interface{}
}

// TextPreProcessor defines the interface for pre-processing text
type TextPreProcessor interface {
	PreProcess(ctx context.Context, text string) (string, error)
}

// PromptGenerator defines the interface for generating prompts
type PromptGenerator interface {
	GeneratePrompt(ctx context.Context, text string) (string, error)
}

// ResponseHandler defines the interface for handling LLM responses
type ResponseHandler interface {
	HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error)
}

// Processor defines the interface for processors
type Processor interface {
	// GetName returns the name of the processor
	GetName() string

	// GetSupportedContentTypes returns the content types this processor can handle
	GetSupportedContentTypes() []string

	// Process processes a ProcessItem
	Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error)

	// ProcessBatch processes a batch of ProcessItems
	ProcessBatch(ctx context.Context, items []*data.ProcessItem) ([]*data.ProcessItem, error)

	// ProcessSource processes all items from a source
	ProcessSource(ctx context.Context, source data.ProcessItemSource, batchSize, workers int) ([]*data.ProcessItem, error)
}

// BaseProcessor provides a base implementation for processors
type BaseProcessor struct {
	name            string
	contentTypes    []string
	llmClient       llm.Client
	preProcessor    TextPreProcessor
	promptGenerator PromptGenerator
	responseHandler ResponseHandler
	options         Options
}

// BaseResponseHandler provides common response handling functionality
type BaseResponseHandler struct {
	ProcessorType    string
	DefaultResponder func() interface{}
	// Fields specifies the fields to extract from the response
	Fields map[string]FieldMapper
	// ResultStruct is a pointer to the struct type to map to (used for automatic mapping)
	ResultStruct interface{}
	// DynamicValidators stores dynamically added validation functions
	DynamicValidators map[string]func(interface{}) interface{}
}

// FieldMapper defines how to map a field from the response to the result
type FieldMapper struct {
	// DefaultValue is the default value to use if the field is missing
	DefaultValue interface{}
	// Transform is an optional function to transform the field value
	Transform func(interface{}) interface{}
}

// CleanResponseString removes markdown code blocks from a response string
func (h *BaseResponseHandler) CleanResponseString(response string) string {
	cleanResponse := response
	// Handle multi-line strings with different code block formats
	if strings.Contains(cleanResponse, "```") {
		// Try to extract content from code blocks with language specifiers
		if strings.Contains(cleanResponse, "```json") {
			// Find all content between ```json and ``` markers
			parts := strings.Split(cleanResponse, "```json")
			if len(parts) > 1 {
				codeContent := parts[1]
				endPos := strings.Index(codeContent, "```")
				if endPos != -1 {
					cleanResponse = strings.TrimSpace(codeContent[:endPos])
				}
			}
		} else {
			// Try to extract content from generic code blocks
			parts := strings.Split(cleanResponse, "```")
			if len(parts) >= 3 { // At least one complete code block
				// Take the content of the first code block
				cleanResponse = strings.TrimSpace(parts[1])
			}
		}
	} else {
		// Handle inline code with backticks
		if strings.HasPrefix(cleanResponse, "`") && strings.HasSuffix(cleanResponse, "`") {
			cleanResponse = cleanResponse[1 : len(cleanResponse)-1]
			cleanResponse = strings.TrimSpace(cleanResponse)
		}
	}
	return cleanResponse
}

// ParseLLMResponse handles common LLM response parsing patterns
func (h *BaseResponseHandler) ParseLLMResponse(responseData interface{}) (map[string]interface{}, bool, interface{}) {
	// Handle string responses from LLM
	if strResponse, ok := responseData.(string); ok {
		// Clean the response by removing markdown
		cleanResponse := h.CleanResponseString(strResponse)

		// Try to parse the string as JSON
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(cleanResponse), &data); err != nil {
			// Before giving up, check if the response itself contains a code block
			// This handles the case where the LLM wrapped the response in markdown but we failed to clean it
			if strings.Contains(strResponse, "```") || strings.Contains(strResponse, "`{") {
				// Try again with a more aggressive markdown cleanup
				jsonStartIndex := strings.Index(strResponse, "{")
				jsonEndIndex := strings.LastIndex(strResponse, "}")

				if jsonStartIndex != -1 && jsonEndIndex != -1 && jsonEndIndex > jsonStartIndex {
					potentialJSON := strResponse[jsonStartIndex : jsonEndIndex+1]
					// Try to parse this potential JSON
					if err := json.Unmarshal([]byte(potentialJSON), &data); err == nil {
						// Successfully parsed JSON within the response
						responseData = data
						return data, true, nil
					}
				}
			}

			// If all parsing attempts fail, return default response
			defaultResponse := h.DefaultResponder()

			// Add the raw response
			result := map[string]interface{}{
				"response":       strResponse,
				"processor_type": h.ProcessorType,
			}

			// Merge the default response fields into result
			if defaultMap, ok := defaultResponse.(map[string]interface{}); ok {
				for k, v := range defaultMap {
					if k != "processor_type" && k != "response" {
						result[k] = v
					}
				}
			}

			return result, false, nil
		}

		// If parsing succeeds, use the parsed data
		responseData = data
	}

	// Convert the response data to map
	data, ok := responseData.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid response data format: %T", responseData)
	}

	// Extract debug info if it exists
	var debugInfo interface{}
	if debug, exists := data["debug"]; exists {
		debugInfo = debug
	}

	// Check if we got a non-JSON response wrapped in a "response" field
	if response, exists := data["response"]; exists && len(data) <= 2 { // data has only response and maybe debug
		// If this response is potentially a string containing JSON, try to parse it
		if responseStr, ok := response.(string); ok && (strings.Contains(responseStr, "{") || strings.Contains(responseStr, "[")) {
			// Clean the response string
			cleanResponse := h.CleanResponseString(responseStr)

			// Try to parse as JSON
			var nestedData map[string]interface{}
			if err := json.Unmarshal([]byte(cleanResponse), &nestedData); err == nil {
				// Successfully parsed JSON from response field
				// Merge the nested data with processor_type
				nestedData["processor_type"] = h.ProcessorType
				return nestedData, true, debugInfo
			}
		}

		// This is a fallback case where the LLM didn't produce valid JSON
		defaultResponse := h.DefaultResponder()

		// Start with the default values
		result := map[string]interface{}{
			"response":       response,
			"processor_type": h.ProcessorType,
		}

		// Merge the default response fields into result
		if defaultMap, ok := defaultResponse.(map[string]interface{}); ok {
			for k, v := range defaultMap {
				if k != "processor_type" && k != "response" {
					result[k] = v
				}
			}
		}

		return result, false, debugInfo
	}

	return data, true, debugInfo
}

// MapResponseToResult maps fields from data to a result map based on field definitions
func (h *BaseResponseHandler) MapResponseToResult(data map[string]interface{}) map[string]interface{} {
	// Start with processor type
	result := map[string]interface{}{
		"processor_type": h.ProcessorType,
	}

	// If we have field definitions, use them to map fields
	if h.Fields != nil {
		for fieldName, mapper := range h.Fields {
			// Extract the value
			value, exists := data[fieldName]

			// Apply dynamic validator if one exists for this field
			if exists && value != nil && h.DynamicValidators != nil {
				if validator, ok := h.DynamicValidators[fieldName]; ok {
					value = validator(value)
				}
			}

			// Apply transformation if provided and value exists
			if exists && mapper.Transform != nil {
				value = mapper.Transform(value)
			}

			// Use default if value doesn't exist
			if !exists || value == nil {
				value = mapper.DefaultValue
			}

			// Add to result
			result[fieldName] = value
		}
	} else {
		// No field definitions, copy all fields except debug and processor_type
		for k, v := range data {
			if k != "debug" && k != "processor_type" {
				// Apply dynamic validator if one exists for this field
				if h.DynamicValidators != nil {
					if validator, ok := h.DynamicValidators[k]; ok {
						v = validator(v)
					}
				}
				result[k] = v
			}
		}
	}

	return result
}

// MapToStruct maps the data to a typed struct using reflection based on json tags
func (h *BaseResponseHandler) MapToStruct(data map[string]interface{}) interface{} {
	if h.ResultStruct == nil {
		// If no result struct is provided, return as map
		return h.MapResponseToResult(data)
	}

	// Get a map with all fields with defaults applied
	resultMap := h.MapResponseToResult(data)

	// Use reflection to map to struct
	result := reflect.New(reflect.TypeOf(h.ResultStruct).Elem()).Interface()

	// Get the struct value and type
	resultValue := reflect.ValueOf(result).Elem()
	resultType := resultValue.Type()

	// Iterate over struct fields
	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)

		// Get the JSON tag
		tag := field.Tag.Get("json")
		if tag == "" {
			// No JSON tag, use field name in lowercase
			tag = strings.ToLower(field.Name)
		} else {
			// Extract the name part of the tag (before any comma)
			tag = strings.Split(tag, ",")[0]
		}

		// Skip if field is processor_type and we're setting it automatically
		if tag == "processor_type" {
			fieldValue := resultValue.Field(i)
			if fieldValue.Kind() == reflect.String && fieldValue.CanSet() {
				fieldValue.SetString(h.ProcessorType)
			}
			continue
		}

		// Get the value from the map
		if mapValue, ok := resultMap[tag]; ok && mapValue != nil {
			fieldValue := resultValue.Field(i)

			// Only set if field is settable
			if !fieldValue.CanSet() {
				continue
			}

			// Convert the value based on field type
			mapped := h.mapValueToField(mapValue, fieldValue)
			if mapped {
				continue
			}

			// If mapping failed and this is a slice, try special handling
			if fieldValue.Kind() == reflect.Slice {
				h.mapSlice(mapValue, fieldValue)
			}
		}
	}

	return result
}

// mapValueToField attempts to map a value to a struct field based on type
func (h *BaseResponseHandler) mapValueToField(value interface{}, field reflect.Value) bool {
	// Get the value as reflect.Value
	valueRefl := reflect.ValueOf(value)

	// Handle nil value
	if value == nil {
		return false
	}

	// If types are directly assignable
	if valueRefl.Type().AssignableTo(field.Type()) {
		field.Set(valueRefl)
		return true
	}

	// Type conversions
	switch field.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			field.SetString(str)
			return true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, ok := value.(float64); ok {
			field.SetInt(int64(num))
			return true
		} else if num, ok := value.(int); ok {
			field.SetInt(int64(num))
			return true
		} else if num, ok := value.(int64); ok {
			field.SetInt(num)
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if num, ok := value.(float64); ok {
			field.SetUint(uint64(num))
			return true
		} else if num, ok := value.(uint); ok {
			field.SetUint(uint64(num))
			return true
		} else if num, ok := value.(uint64); ok {
			field.SetUint(num)
			return true
		}
	case reflect.Float32, reflect.Float64:
		if num, ok := value.(float64); ok {
			field.SetFloat(num)
			return true
		} else if num, ok := value.(float32); ok {
			field.SetFloat(float64(num))
			return true
		} else if num, ok := value.(int); ok {
			field.SetFloat(float64(num))
			return true
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			field.SetBool(b)
			return true
		}
	}

	return false
}

// mapSlice handles mapping to slice types
func (h *BaseResponseHandler) mapSlice(value interface{}, field reflect.Value) {
	// Handle []string specifically, which is common
	if field.Type().Elem().Kind() == reflect.String {
		// Try to convert various slice types to []string
		switch v := value.(type) {
		case []string:
			// Direct assignment
			field.Set(reflect.ValueOf(v))

		case []interface{}:
			// Convert []interface{} to []string
			strSlice := make([]string, 0, len(v))
			for _, item := range v {
				if s, ok := item.(string); ok {
					strSlice = append(strSlice, s)
				}
			}

			// Create a new slice value and assign
			newSlice := reflect.MakeSlice(field.Type(), len(strSlice), len(strSlice))
			for i, s := range strSlice {
				newSlice.Index(i).SetString(s)
			}
			field.Set(newSlice)
		}
	}
}

// AutoProcessResponse is a complete response processing workflow that handles:
// - Parsing the LLM response
// - Mapping to a result struct
// - Handling debug info
// This reduces boilerplate code in individual processors.
func (h *BaseResponseHandler) AutoProcessResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	// Parse the LLM response
	data, validJSON, debugInfo := h.ParseLLMResponse(responseData)
	if data == nil {
		return nil, fmt.Errorf("failed to parse response data")
	}

	// If we don't have valid JSON structure, return the default response
	if !validJSON {
		return data, nil
	}

	// Map fields from data to result struct using reflection
	result := h.MapToStruct(data)

	// Add debug info if needed
	if debugInfo != nil {
		// If the result is already a map, just add debug info
		if resultMap, ok := result.(map[string]interface{}); ok {
			resultMap["debug"] = debugInfo
			return resultMap, nil
		}

		// Otherwise, convert the struct to a map to include debug info
		resultMap := map[string]interface{}{
			"processor_type": h.ProcessorType,
			"debug":          debugInfo,
		}

		// Use reflection to copy fields from the struct to the map
		resultValue := reflect.ValueOf(result).Elem()
		resultType := resultValue.Type()

		for i := 0; i < resultType.NumField(); i++ {
			field := resultType.Field(i)
			fieldValue := resultValue.Field(i)

			// Get the JSON tag name
			tag := field.Tag.Get("json")
			if tag == "" {
				// No JSON tag, use field name in lowercase
				tag = strings.ToLower(field.Name)
			} else {
				// Extract the name part of the tag (before any comma)
				tag = strings.Split(tag, ",")[0]
			}

			// Skip omitempty fields that are empty
			if strings.Contains(field.Tag.Get("json"), "omitempty") {
				// Skip empty slices
				if fieldValue.Kind() == reflect.Slice && fieldValue.Len() == 0 {
					continue
				}

				// Skip empty strings
				if fieldValue.Kind() == reflect.String && fieldValue.String() == "" {
					continue
				}
			}

			// Add field to map
			resultMap[tag] = fieldValue.Interface()
		}

		return resultMap, nil
	}

	return result, nil
}

// NewBaseProcessor creates a new base processor
func NewBaseProcessor(name string, contentTypes []string, llmClient llm.Client,
	preProcessor TextPreProcessor, promptGenerator PromptGenerator,
	responseHandler ResponseHandler, options Options) *BaseProcessor {

	if contentTypes == nil {
		contentTypes = []string{"text"}
	}

	return &BaseProcessor{
		name:            name,
		contentTypes:    contentTypes,
		llmClient:       llmClient,
		preProcessor:    preProcessor,
		promptGenerator: promptGenerator,
		responseHandler: responseHandler,
		options:         options,
	}
}

// GetName returns the processor name
func (p *BaseProcessor) GetName() string {
	return p.name
}

// GetSupportedContentTypes returns content types this processor can handle
func (p *BaseProcessor) GetSupportedContentTypes() []string {
	return p.contentTypes
}

// Process processes a ProcessItem
func (p *BaseProcessor) Process(ctx context.Context, item *data.ProcessItem) (*data.ProcessItem, error) {
	// Validate content type
	contentTypeSupported := false
	for _, ct := range p.contentTypes {
		if ct == item.ContentType {
			contentTypeSupported = true
			break
		}
	}

	if !contentTypeSupported {
		return nil, fmt.Errorf("unsupported content type: %s", item.ContentType)
	}

	// Clone the item to avoid modifying the original
	result, err := item.Clone()
	if err != nil {
		return nil, err
	}

	// Get text content based on the content type
	var textContent string

	if item.ContentType == "text" {
		// Get text content directly
		textContent, err = item.GetTextContent()
		if err != nil {
			return nil, err
		}
	} else if item.ContentType == "json" {
		// For JSON content, either:
		// 1. Use "text" field if available in the JSON
		// 2. Use "response" field if available
		// 3. Or convert the entire JSON to text as fallback
		jsonContent, ok := item.Content.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid JSON content format")
		}

		// Try to extract text from the JSON
		if text, ok := jsonContent["text"].(string); ok {
			textContent = text
		} else if text, ok := jsonContent["response"].(string); ok {
			textContent = text
		} else if originalText, ok := item.Metadata["original_text"].(string); ok {
			// Try to get original text from metadata if available
			textContent = originalText
		} else {
			// Use the first text field we can find
			foundText := false
			for _, value := range jsonContent {
				if text, ok := value.(string); ok {
					textContent = text
					foundText = true
					break
				}
			}

			// If we still don't have text, convert the JSON to string
			if !foundText {
				jsonBytes, err := json.Marshal(jsonContent)
				if err != nil {
					return nil, fmt.Errorf("failed to convert JSON to text: %w", err)
				}
				textContent = string(jsonBytes)
			}
		}
	}

	// Run LLM processing if available
	if p.llmClient != nil {
		// Check if debug is enabled in options
		debugEnabled := false
		if p.options.LLMOptions != nil {
			if debug, ok := p.options.LLMOptions["debug"].(bool); ok {
				debugEnabled = debug
			}
		}

		// Pre-process if needed
		if p.preProcessor != nil {
			textContent, err = p.preProcessor.PreProcess(ctx, textContent)
			if err != nil {
				return nil, err
			}
		}

		// Generate prompt if needed
		prompt := textContent
		if p.promptGenerator != nil {
			prompt, err = p.promptGenerator.GeneratePrompt(ctx, textContent)
			if err != nil {
				return nil, err
			}
		}

		// Print debug information if enabled
		if debugEnabled {
			DebugLLMInteraction(prompt, "") // Print the prompt before calling LLM
		}

		// Call LLM
		llmResponse, err := p.llmClient.Complete(ctx, prompt, p.options.LLMOptions)
		if err != nil {
			return nil, err
		}

		// Print debug information if enabled
		if debugEnabled {
			DebugLLMInteraction(prompt, llmResponse) // Print full interaction
		}

		// Store debug info in a map if debug is enabled
		var debugInfo map[string]interface{}
		if debugEnabled {
			debugInfo = map[string]interface{}{
				"prompt":       prompt,
				"raw_response": llmResponse,
			}
		}

		// Handle response
		if p.responseHandler != nil {
			processedContent, err := p.responseHandler.HandleResponse(ctx, textContent, llmResponse)
			if err != nil {
				return nil, err
			}

			// Add debug info to processed content if available
			if debugEnabled && debugInfo != nil {
				// If the result is a map, add debug info directly
				if contentMap, ok := processedContent.(map[string]interface{}); ok {
					contentMap["debug"] = debugInfo
					processedContent = contentMap
				} else {
					// For struct responses, we'll handle debug in a different way below
				}
			}

			// Update the content with the processed result
			result.Content = processedContent

			// If content is a string, keep content type as text
			// otherwise change to the appropriate type
			if _, ok := processedContent.(string); !ok {
				result.ContentType = "json"
			} else {
				result.ContentType = "text"
			}

			// Add processing info, checking if processor_type already exists in the response
			if contentMap, ok := processedContent.(map[string]interface{}); ok && contentMap["processor_type"] != nil {
				// Use the processor_type from the response
				result.AddProcessingInfo(p.name, processedContent)
			} else {
				// For struct responses, convert to map first
				// This handles cases like SentimentResult, IntentResult, etc.
				if reflect.TypeOf(processedContent) != nil && reflect.TypeOf(processedContent).Kind() == reflect.Ptr {
					// Use reflection to convert struct to map
					val := reflect.ValueOf(processedContent).Elem()
					if val.Kind() == reflect.Struct {
						structMap := make(map[string]interface{})
						structType := val.Type()

						// First see if struct has a ProcessorType field
						var hasProcessorType bool
						var processorTypeValue string

						// Check each field in the struct
						for i := 0; i < val.NumField(); i++ {
							field := structType.Field(i)

							// Get the field's JSON tag
							tag := field.Tag.Get("json")
							if tag == "" {
								tag = strings.ToLower(field.Name)
							} else {
								tag = strings.Split(tag, ",")[0]
							}

							// Skip if the tag is "-" (meaning don't include in JSON)
							if tag == "-" {
								continue
							}

							// Get the field value
							fieldValue := val.Field(i).Interface()
							structMap[tag] = fieldValue

							// Check if this is the processor_type field
							if tag == "processor_type" {
								hasProcessorType = true
								if strValue, ok := fieldValue.(string); ok {
									processorTypeValue = strValue
								}
							}
						}

						// Add debug info to the struct map if enabled
						if debugEnabled && debugInfo != nil {
							structMap["debug"] = debugInfo
						}

						// If the struct has a processor_type, use it
						if hasProcessorType && processorTypeValue != "" {
							result.AddProcessingInfo(p.name, structMap)
							result.Content = processedContent // Keep the original content
						} else {
							// Add the processor type to the map
							structMap["processor_type"] = p.name
							result.AddProcessingInfo(p.name, structMap)
						}

						return result, nil
					}
				}

				// If not a struct or conversion failed, use the default processor_type
				processingInfo := map[string]interface{}{
					"processor_type": p.name,
				}

				// Add debug info if enabled
				if debugEnabled && debugInfo != nil {
					processingInfo["debug"] = debugInfo
				}

				result.AddProcessingInfo(p.name, processingInfo)
			}
		} else {
			// Default behavior: replace content with LLM response
			result.Content = llmResponse

			// If response is a string, assume it's text
			if _, ok := llmResponse.(string); ok {
				result.ContentType = "text"
			} else {
				result.ContentType = "json"
			}

			// Add processing info with the proper processor type for non-LLM processing
			processingInfo := map[string]interface{}{
				"processor_type": p.name,
			}

			// Add debug info if enabled
			if debugEnabled && debugInfo != nil {
				processingInfo["debug"] = debugInfo
			}

			result.AddProcessingInfo(p.name, processingInfo)
		}
	} else {
		// Add processing info with the proper processor type for non-LLM processing
		result.AddProcessingInfo(p.name, map[string]interface{}{
			"processor_type": p.name,
		})
	}

	// Store original text in metadata if not already present
	if _, exists := result.Metadata["original_text"]; !exists {
		if result.Metadata == nil {
			result.Metadata = make(map[string]interface{})
		}
		result.Metadata["original_text"] = textContent
	}

	return result, nil
}

// ProcessBatch processes a batch of items
func (p *BaseProcessor) ProcessBatch(ctx context.Context, items []*data.ProcessItem) ([]*data.ProcessItem, error) {
	results := make([]*data.ProcessItem, len(items))

	for i, item := range items {
		result, err := p.Process(ctx, item)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// ProcessSource processes all items from a source
func (p *BaseProcessor) ProcessSource(ctx context.Context, source data.ProcessItemSource, batchSize, workers int) ([]*data.ProcessItem, error) {
	processor := data.NewProcessItemParallelProcessor(source, batchSize, workers)
	defer processor.Close()

	return processor.ProcessAll(ctx, p.Process)
}

// InitializeBuiltInProcessors ensures all built-in processors are registered before use
func InitializeBuiltInProcessors() {
	// This function must be called early in the application to ensure
	// all processor init() functions have run and registered their processors

	// Force loading of processor packages by name
	// The side effect of importing these packages is that their init() functions will run
	// and register themselves with the processor registry

	// We don't need any actual code here, just the import side effects
}

// init runs automatically and calls RegisterBuiltInProcessors to ensure processors are registered
func init() {
	// Make sure our init() function runs after all processors are registered
	// This happens automatically due to Go's package initialization order
}

// NewResponseHandler creates a BaseResponseHandler from a result struct
// This automatically configures the handler based on the struct's fields and tags
func NewResponseHandler(processorType string, resultStruct interface{}, customDefaults ...map[string]interface{}) *BaseResponseHandler {
	// Create a default response handler
	handler := &BaseResponseHandler{
		ProcessorType: processorType,
		ResultStruct:  resultStruct,
		Fields:        make(map[string]FieldMapper),
	}

	// Set the default responder based on the result struct
	handler.DefaultResponder = func() interface{} {
		return handler.createDefaultResponse()
	}

	// Configure fields based on the result struct
	handler.configureFieldsFromStruct()

	// Apply processor-specific defaults
	handler.applyProcessorDefaults()

	// Apply any custom defaults (these override the processor-specific ones)
	if len(customDefaults) > 0 && customDefaults[0] != nil {
		for field, value := range customDefaults[0] {
			if _, exists := handler.Fields[field]; exists {
				handler.Fields[field] = FieldMapper{
					DefaultValue: value,
					// Keep any existing transform function
					Transform: handler.Fields[field].Transform,
				}
			}
		}
	}

	return handler
}

// applyProcessorDefaults applies default values specific to each processor type
func (h *BaseResponseHandler) applyProcessorDefaults() {
	// Use reflection to get processor-specific default values from the struct itself
	if h.ResultStruct != nil {
		structValue := reflect.ValueOf(h.ResultStruct).Elem()
		structType := structValue.Type()

		// Look for methods on the struct to provide defaults
		// E.g., a method named "DefaultValues() map[string]interface{}"
		defaultsMethod := reflect.ValueOf(h.ResultStruct).MethodByName("DefaultValues")
		if defaultsMethod.IsValid() {
			// Call the DefaultValues method to get custom defaults
			results := defaultsMethod.Call(nil)
			if len(results) > 0 {
				if defaults, ok := results[0].Interface().(map[string]interface{}); ok {
					// Apply custom defaults from the method
					for field, value := range defaults {
						h.updateFieldMapper(field, value, nil)
					}
					return
				}
			}
		}

		// No custom defaults method, use reflection to scan for tags
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)

			// Check for a "default" tag
			defaultTag := field.Tag.Get("default")
			if defaultTag != "" {
				// Get the field's JSON name
				jsonTag := field.Tag.Get("json")
				fieldName := strings.ToLower(field.Name)
				if jsonTag != "" {
					// Extract the name part of the tag (before any comma)
					fieldName = strings.Split(jsonTag, ",")[0]
				}

				// Convert the default value based on field type
				var defaultValue interface{}
				switch field.Type.Kind() {
				case reflect.String:
					defaultValue = defaultTag
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if val, err := strconv.ParseInt(defaultTag, 10, 64); err == nil {
						defaultValue = val
					}
				case reflect.Float32, reflect.Float64:
					if val, err := strconv.ParseFloat(defaultTag, 64); err == nil {
						defaultValue = val
					}
				case reflect.Bool:
					if val, err := strconv.ParseBool(defaultTag); err == nil {
						defaultValue = val
					}
				}

				if defaultValue != nil {
					h.updateFieldMapper(fieldName, defaultValue, nil)
				}
			}
		}

		// Special handling for complex field types that can't be handled by tags
		// Look for validator/transform methods for specific fields
		// E.g. a method named "ValidateAttributes() interface{}"
		for fieldName := range h.Fields {
			// Build the validator method name: "Validate" + Title case field name
			methodName := "Validate" + strings.Title(fieldName)
			validatorMethod := reflect.ValueOf(h.ResultStruct).MethodByName(methodName)

			if validatorMethod.IsValid() {
				// Call the validator method to get transform function
				results := validatorMethod.Call(nil)
				if len(results) > 0 {
					if transformFn, ok := results[0].Interface().(func(interface{}) interface{}); ok {
						// Update the field mapper with this transform function
						defaultValue := h.Fields[fieldName].DefaultValue
						h.Fields[fieldName] = FieldMapper{
							DefaultValue: defaultValue,
							Transform:    transformFn,
						}
					}
				}
			}
		}
	}
}

// updateFieldMapper updates a field mapper with a new default value and optional transform
func (h *BaseResponseHandler) updateFieldMapper(field string, defaultValue interface{}, transform func(interface{}) interface{}) {
	if _, exists := h.Fields[field]; exists {
		// Keep existing transform if a new one isn't provided
		existingTransform := h.Fields[field].Transform
		if transform == nil {
			transform = existingTransform
		}

		h.Fields[field] = FieldMapper{
			DefaultValue: defaultValue,
			Transform:    transform,
		}
	}
}

// createDefaultResponse creates a default response map based on the result struct
func (h *BaseResponseHandler) createDefaultResponse() map[string]interface{} {
	if h.ResultStruct == nil {
		return map[string]interface{}{"processor_type": h.ProcessorType}
	}

	// Check if the struct has a DefaultValues method
	defaultsMethod := reflect.ValueOf(h.ResultStruct).MethodByName("DefaultValues")
	if defaultsMethod.IsValid() {
		// Call the DefaultValues method to get custom defaults
		results := defaultsMethod.Call(nil)
		if len(results) > 0 {
			if defaults, ok := results[0].Interface().(map[string]interface{}); ok {
				// Add processor_type
				defaults["processor_type"] = h.ProcessorType
				return defaults
			}
		}
	}

	// Fallback to using the general GetDefaultValues function
	defaults := GetDefaultValues(h.ResultStruct)
	defaults["processor_type"] = h.ProcessorType
	return defaults
}

// configureFieldsFromStruct automatically configures field mappers from the result struct
func (h *BaseResponseHandler) configureFieldsFromStruct() {
	if h.ResultStruct == nil || h.Fields == nil {
		return
	}

	// Get the struct value and type
	structValue := reflect.ValueOf(h.ResultStruct).Elem()
	structType := structValue.Type()

	// Iterate over struct fields
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		// Get the JSON tag name
		tag := field.Tag.Get("json")
		if tag == "" {
			// No JSON tag, use field name in lowercase
			tag = strings.ToLower(field.Name)
		} else {
			// Extract the name part of the tag (before any comma)
			tag = strings.Split(tag, ",")[0]
		}

		// Skip processor_type as it's handled automatically
		if tag == "processor_type" {
			continue
		}

		// Add a field mapper with appropriate default value
		var defaultValue interface{}

		switch field.Type.Kind() {
		case reflect.String:
			defaultValue = ""
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			defaultValue = int64(0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			defaultValue = uint64(0)
		case reflect.Float32, reflect.Float64:
			defaultValue = float64(0)
		case reflect.Bool:
			defaultValue = false
		case reflect.Slice:
			// Create an empty slice of the appropriate type
			defaultValue = reflect.MakeSlice(field.Type, 0, 0).Interface()
		case reflect.Map:
			// Create an empty map of the appropriate type
			defaultValue = reflect.MakeMap(field.Type).Interface()
		case reflect.Struct:
			// For embedded structs, use a nil pointer to indicate empty
			defaultValue = nil
		default:
			defaultValue = nil
		}

		// Add to fields map
		h.Fields[tag] = FieldMapper{
			DefaultValue: defaultValue,
		}

		// Add special handling for common types
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			// Add transform for []string type fields
			h.Fields[tag] = FieldMapper{
				DefaultValue: []string{},
				Transform: func(val interface{}) interface{} {
					// Handle both []string and []interface{} types
					if strSlice, ok := val.([]string); ok {
						return strSlice
					}

					if items, ok := val.([]interface{}); ok {
						result := make([]string, 0, len(items))
						for _, item := range items {
							if s, ok := item.(string); ok {
								result = append(result, s)
							}
						}
						return result
					}

					return []string{}
				},
			}
		}
	}
}

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

// Helper functions for processors

// GetStringValue safely gets a string value from an interface map
func GetStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// GetFloatValue safely gets a float value from an interface map
func GetFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

// GetIntValue safely gets an int value from an interface map
func GetIntValue(m map[string]interface{}, key string) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	if val, ok := m[key].(int); ok {
		return val
	}
	return 0
}

// GetBoolValue safely gets a bool value from an interface map
func GetBoolValue(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}

// HandleResponse implements ResponseHandler interface
func (h *BaseResponseHandler) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	return h.AutoProcessResponse(ctx, text, responseData)
}
