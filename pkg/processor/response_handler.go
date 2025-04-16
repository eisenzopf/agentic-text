package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// FieldMapper defines how to map a field from the response to the result
type FieldMapper struct {
	// DefaultValue is the default value to use if the field is missing
	DefaultValue interface{}
	// Transform is an optional function to transform the field value
	Transform func(interface{}) interface{}
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
	// validateStructure determines if strict structural validation should be performed
	validateStructure bool
}

// CleanResponseString removes markdown code blocks from a response string
func (h *BaseResponseHandler) CleanResponseString(response string) string {
	cleanResponse := response

	// Handle multi-line strings with different code block formats
	if strings.Contains(cleanResponse, "```") {
		// Try multiple strategies to extract code blocks

		// Strategy 1: Extract from ```json blocks
		if strings.Contains(cleanResponse, "```json") {
			parts := strings.Split(cleanResponse, "```json")
			if len(parts) > 1 {
				codeContent := parts[1]
				endPos := strings.Index(codeContent, "```")
				if endPos != -1 {
					cleanResponse = strings.TrimSpace(codeContent[:endPos])
					return cleanResponse
				}
			}
		}

		// Strategy 2: Look for any language specifier
		languageSpecifiers := []string{"```javascript", "```js", "```python", "```go", "```java", "```typescript", "```ts"}
		for _, specifier := range languageSpecifiers {
			if strings.Contains(cleanResponse, specifier) {
				parts := strings.Split(cleanResponse, specifier)
				if len(parts) > 1 {
					codeContent := parts[1]
					endPos := strings.Index(codeContent, "```")
					if endPos != -1 {
						cleanResponse = strings.TrimSpace(codeContent[:endPos])
						return cleanResponse
					}
				}
			}
		}

		// Strategy 3: Extract content from generic code blocks
		parts := strings.Split(cleanResponse, "```")
		if len(parts) >= 3 { // At least one complete code block
			// Take the content of the first code block
			cleanResponse = strings.TrimSpace(parts[1])
			return cleanResponse
		}

		// Strategy 4: If all else fails, try to find JSON between the first { and last }
		jsonStartIndex := strings.Index(cleanResponse, "{")
		jsonEndIndex := strings.LastIndex(cleanResponse, "}")

		if jsonStartIndex != -1 && jsonEndIndex != -1 && jsonEndIndex > jsonStartIndex {
			potentialJSON := cleanResponse[jsonStartIndex : jsonEndIndex+1]
			// Verify it's valid JSON
			var testJSON map[string]interface{}
			if err := json.Unmarshal([]byte(potentialJSON), &testJSON); err == nil {
				return potentialJSON
			}
		}
	} else {
		// Handle inline code with backticks
		if strings.HasPrefix(cleanResponse, "`") && strings.HasSuffix(cleanResponse, "`") {
			cleanResponse = cleanResponse[1 : len(cleanResponse)-1]
			cleanResponse = strings.TrimSpace(cleanResponse)
		}

		// Handle potential JSON without code blocks
		jsonStartIndex := strings.Index(cleanResponse, "{")
		jsonEndIndex := strings.LastIndex(cleanResponse, "}")

		if jsonStartIndex != -1 && jsonEndIndex != -1 && jsonEndIndex > jsonStartIndex {
			potentialJSON := cleanResponse[jsonStartIndex : jsonEndIndex+1]
			// Verify it's valid JSON
			var testJSON map[string]interface{}
			if err := json.Unmarshal([]byte(potentialJSON), &testJSON); err == nil {
				return potentialJSON
			}
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
	} else if field.Type().Elem().Kind() == reflect.Struct {
		// Handle slice of structs (like []AttributeDefinition)
		if items, ok := value.([]interface{}); ok && len(items) > 0 {
			// Create a new slice of the appropriate struct type
			elemType := field.Type().Elem()
			newSlice := reflect.MakeSlice(field.Type(), 0, len(items))

			// Process each item in the array
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					// Create a new struct of the appropriate type
					newStruct := reflect.New(elemType).Elem()

					// Map fields from the map to the struct
					for i := 0; i < elemType.NumField(); i++ {
						structField := elemType.Field(i)

						// Get JSON tag name
						tag := structField.Tag.Get("json")
						if tag == "" {
							tag = strings.ToLower(structField.Name)
						} else {
							tag = strings.Split(tag, ",")[0]
						}

						// Find the value in the map
						if mapValue, exists := itemMap[tag]; exists && mapValue != nil {
							fieldValue := newStruct.Field(i)

							// Only set if field is settable
							if !fieldValue.CanSet() {
								continue
							}

							// Convert and set the value
							if h.mapValueToField(mapValue, fieldValue) {
								// Value was set successfully
								continue
							}

							// If basic mapping failed and the target is a slice, try mapSlice
							if fieldValue.Kind() == reflect.Slice {
								h.mapSlice(mapValue, fieldValue)
							}

							// For complex types, might need additional handling here
						}
					}

					// Add the new struct to the slice
					newSlice = reflect.Append(newSlice, newStruct)
				}
			}

			// Only set the field if we have elements
			if newSlice.Len() > 0 {
				field.Set(newSlice)
			}
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
		return data, nil // data here contains the non-JSON response and defaults
	}

	// --- Structural Validation Step ---
	if h.validateStructure {
		// Attempt to map the data to the struct to check structural compatibility.
		// MapToStruct internally uses MapResponseToResult which applies defaults and validators.
		tentativeResult := h.MapToStruct(data)

		// Simple check: If mapping resulted in nil or didn't produce the expected type,
		// consider it a structural validation failure.
		// More sophisticated checks could be added here (e.g., check required fields).
		if tentativeResult == nil || reflect.TypeOf(tentativeResult) != reflect.TypeOf(h.ResultStruct) {
			// Validation failed, return the default response object.
			// We need to ensure the default response includes the processor_type.
			defaultResponseMap := h.createDefaultResponse()
			// Add debug info to the default response if available
			if debugInfo != nil {
				defaultResponseMap["debug"] = debugInfo
			}
			return defaultResponseMap, nil
		}
		// If validation passes, we can proceed with the result from the tentative mapping.
		result := tentativeResult

		// Add debug info if needed (handling map vs struct)
		if debugInfo != nil {
			AddDebugInfoToResult(&result, debugInfo, h.ProcessorType)
		}
		return result, nil

	} else {
		// --- No Structural Validation ---
		// Proceed with mapping without the strict structural check
		result := h.MapToStruct(data)

		// Add debug info if needed (handling map vs struct)
		if debugInfo != nil {
			AddDebugInfoToResult(&result, debugInfo, h.ProcessorType)
		}
		return result, nil
	}
}

// HandleResponse implements ResponseHandler interface
func (h *BaseResponseHandler) HandleResponse(ctx context.Context, text string, responseData interface{}) (interface{}, error) {
	return h.AutoProcessResponse(ctx, text, responseData)
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
