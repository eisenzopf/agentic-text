package processor

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// GenerateJSONExample generates a sample JSON structure from a struct
// This is useful for creating example JSON in LLM prompts
func GenerateJSONExample(structType interface{}) string {
	// Create a sample instance of the struct
	val := reflect.ValueOf(structType).Elem()
	sampleStruct := reflect.New(val.Type()).Interface()

	// Use reflection to populate the struct with sample values
	populateSampleValues(reflect.ValueOf(sampleStruct).Elem())

	// Marshal to a map first so we can exclude processor_type
	jsonBytes, err := json.Marshal(sampleStruct)
	if err != nil {
		return "{}"
	}

	// Unmarshal to a map
	var resultMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &resultMap); err != nil {
		return "{}"
	}

	// Remove the processor_type field
	delete(resultMap, "processor_type")

	// Marshal back to JSON with indentation
	prettyBytes, err := json.MarshalIndent(resultMap, "", "  ")
	if err != nil {
		return "{}"
	}

	return string(prettyBytes)
}

// populateSampleValues populates a struct value with sample data
func populateSampleValues(value reflect.Value) {
	// Get the struct type
	typ := value.Type()

	// Iterate through fields
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Check if there's a json tag
		tag := fieldType.Tag.Get("json")
		if tag == "-" {
			continue // Skip fields with json:"-"
		}

		// Get field name from JSON tag or struct field name
		fieldName := fieldType.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		// Generate sample value based on field type
		switch field.Kind() {
		case reflect.String:
			// Look for example comments in the field
			comment := fieldType.Tag.Get("comment")
			if comment != "" {
				field.SetString(comment)
			} else {
				// Use field name as sample
				field.SetString("Example " + fieldName)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(42)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			field.SetUint(42)
		case reflect.Float32, reflect.Float64:
			field.SetFloat(42.5)
		case reflect.Bool:
			field.SetBool(true)
		case reflect.Slice:
			// Create a sample slice with one element
			sliceType := field.Type().Elem()
			sampleSlice := reflect.MakeSlice(field.Type(), 1, 1)

			// If it's a slice of struct, populate the struct
			if sliceType.Kind() == reflect.Struct {
				populateSampleValues(sampleSlice.Index(0))
			} else if sliceType.Kind() == reflect.String {
				sampleSlice.Index(0).SetString("Sample string in array")
			} else if sliceType.Kind() == reflect.Int {
				sampleSlice.Index(0).SetInt(42)
			}

			field.Set(sampleSlice)
		case reflect.Struct:
			populateSampleValues(field)
		case reflect.Ptr:
			// Create a new instance of the pointed-to type and set it
			if field.IsNil() {
				newVal := reflect.New(field.Type().Elem())
				field.Set(newVal)
				populateSampleValues(newVal.Elem())
			}
		}
	}
}

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
