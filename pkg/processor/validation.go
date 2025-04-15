package processor

import (
	"reflect"
	"strconv"
	"strings"
)

// ValidateData is a generic validation function that ensures data returned from LLM responses
// is properly structured before being used in the application
func ValidateData(fieldName string, defaultValue interface{}) func(interface{}) interface{} {
	return func(val interface{}) interface{} {
		// Handle different ways the LLM might return data
		switch v := val.(type) {
		case []interface{}:
			// Direct array format
			if len(v) == 0 {
				return defaultValue
			}

			// Validate items in the array if they're maps
			validItems := make([]interface{}, 0, len(v))
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					// Ensure required field exists
					if GetStringValue(itemMap, fieldName) != "" {
						validItems = append(validItems, itemMap)
					}
				}
			}

			if len(validItems) == 0 {
				return defaultValue
			}
			return validItems

		case map[string]interface{}:
			// Check if data is in a nested field
			if nestedData, ok := v[fieldName].([]interface{}); ok {
				return ValidateData(fieldName, defaultValue)(nestedData)
			}

			// If field name exists directly in the map
			if GetStringValue(v, fieldName) != "" {
				return v
			}

			return defaultValue

		default:
			return defaultValue
		}
	}
}

// DefaultsFromStruct automatically generates default values from a struct using `default` tags
// This simplifies processor definition by extracting defaults from the struct definition
func DefaultsFromStruct(structPtr interface{}) map[string]interface{} {
	defaults := make(map[string]interface{})

	v := reflect.ValueOf(structPtr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip the ProcessorType field which is set by the framework
		if field.Name == "ProcessorType" {
			continue
		}

		// Get JSON field name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Extract field name from json tag (handling cases like "field_name,omitempty")
		jsonName := strings.Split(jsonTag, ",")[0]

		// Get default value from tag
		defaultTag := field.Tag.Get("default")
		if defaultTag != "" {
			// Convert string default value to appropriate type
			switch field.Type.Kind() {
			case reflect.String:
				defaults[jsonName] = defaultTag
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if val, err := parseInt(defaultTag); err == nil {
					defaults[jsonName] = val
				}
			case reflect.Float32, reflect.Float64:
				if val, err := parseFloat(defaultTag); err == nil {
					defaults[jsonName] = val
				}
			case reflect.Bool:
				if val, err := parseBool(defaultTag); err == nil {
					defaults[jsonName] = val
				}
			case reflect.Slice:
				// For slices, default to empty slice if "omitempty" is in the JSON tag
				if strings.Contains(jsonTag, "omitempty") {
					defaults[jsonName] = reflect.MakeSlice(field.Type, 0, 0).Interface()
				}
			}
		} else if strings.Contains(jsonTag, "omitempty") {
			// For fields with omitempty and no explicit default, set appropriate zero value
			switch field.Type.Kind() {
			case reflect.String:
				defaults[jsonName] = ""
			case reflect.Slice:
				defaults[jsonName] = reflect.MakeSlice(field.Type, 0, 0).Interface()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				defaults[jsonName] = 0
			case reflect.Float32, reflect.Float64:
				defaults[jsonName] = 0.0
			case reflect.Bool:
				defaults[jsonName] = false
			}
		}
	}

	return defaults
}

// GetDefaultValues returns default values for any result struct by extracting from struct tags
// This eliminates the need for each processor to implement its own DefaultValues() method
func GetDefaultValues(resultStruct interface{}) map[string]interface{} {
	return DefaultsFromStruct(resultStruct)
}

// Helper functions to parse default values
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
