package processor

import (
	"reflect"
	"strings"
)

// MapResponseToResult maps fields from data to a result map based on field definitions
func MapResponseToResult(data map[string]interface{}, processorType string,
	fields map[string]FieldMapper, dynamicValidators map[string]func(interface{}) interface{}) map[string]interface{} {

	// Start with processor type
	result := map[string]interface{}{
		"processor_type": processorType,
	}

	// If we have field definitions, use them to map fields
	if fields != nil {
		for fieldName, mapper := range fields {
			// Extract the value
			value, exists := data[fieldName]

			// Apply dynamic validator if one exists for this field
			if exists && value != nil && dynamicValidators != nil {
				if validator, ok := dynamicValidators[fieldName]; ok {
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
				if dynamicValidators != nil {
					if validator, ok := dynamicValidators[k]; ok {
						v = validator(v)
					}
				}
				result[k] = v
			}
		}
	}

	return result
}

// MapValueToField attempts to map a value to a struct field based on type
func MapValueToField(value interface{}, field reflect.Value) bool {
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

// MapSlice handles mapping to slice types
func MapSlice(value interface{}, field reflect.Value, mapValueFunc func(interface{}, reflect.Value) bool) {
	// Use passed function or default to MapValueToField if nil
	var mapValueFn func(interface{}, reflect.Value) bool
	if mapValueFunc != nil {
		mapValueFn = mapValueFunc
	} else {
		mapValueFn = MapValueToField
	}

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
							if mapValueFn(mapValue, fieldValue) {
								// Value was set successfully
								continue
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
