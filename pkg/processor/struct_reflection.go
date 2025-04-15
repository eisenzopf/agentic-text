package processor

import (
	"reflect"
	"strings"
)

// ConfigureFieldsFromStruct automatically configures field mappers from a result struct
func ConfigureFieldsFromStruct(resultStruct interface{}, fields map[string]FieldMapper) {
	if resultStruct == nil || fields == nil {
		return
	}

	// Get the struct value and type
	structValue := reflect.ValueOf(resultStruct).Elem()
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
		fields[tag] = FieldMapper{
			DefaultValue: defaultValue,
		}

		// Add special handling for common types
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			// Add transform for []string type fields
			fields[tag] = FieldMapper{
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

// MapToStruct maps data to a typed struct using reflection based on json tags
func MapToStruct(data map[string]interface{}, resultStruct interface{}, processorType string,
	fields map[string]FieldMapper, dynamicValidators map[string]func(interface{}) interface{}) interface{} {

	if resultStruct == nil {
		// If no result struct is provided, return as map
		return MapResponseToResult(data, processorType, fields, dynamicValidators)
	}

	// Get a map with all fields with defaults applied
	resultMap := MapResponseToResult(data, processorType, fields, dynamicValidators)

	// Use reflection to map to struct
	result := reflect.New(reflect.TypeOf(resultStruct).Elem()).Interface()

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
				fieldValue.SetString(processorType)
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
			mapped := MapValueToField(mapValue, fieldValue)
			if mapped {
				continue
			}

			// If mapping failed and this is a slice, try special handling
			if fieldValue.Kind() == reflect.Slice {
				MapSlice(mapValue, fieldValue, nil)
			}
		}
	}

	return result
}

// AddDebugInfoToResult adds debug info to a response result
func AddDebugInfoToResult(result interface{}, debugInfo interface{}, processorType string) interface{} {
	// If there's no debug info, just return the result
	if debugInfo == nil {
		return result
	}

	// If the result is already a map, just add debug info
	if resultMap, ok := result.(map[string]interface{}); ok {
		resultMap["debug"] = debugInfo
		return resultMap
	}

	// Otherwise, convert the struct to a map to include debug info
	resultMap := map[string]interface{}{
		"processor_type": processorType,
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

	return resultMap
}
