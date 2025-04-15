package processor

import (
	"encoding/json"
	"fmt"
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

	// Use reflection to populate the struct with sample values based on default tags
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

		// Check for default tag value
		defaultValue := fieldType.Tag.Get("default")

		// Generate sample value based on field type
		switch field.Kind() {
		case reflect.String:
			if defaultValue != "" {
				// Use the default value if provided
				field.SetString(defaultValue)
			} else {
				// Look for example comments in the field
				comment := fieldType.Tag.Get("comment")
				if comment != "" {
					field.SetString(comment)
				} else {
					// Use field name as sample
					field.SetString("Example " + fieldName)
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if defaultValue != "" {
				// Convert default string to int
				if intVal, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
					field.SetInt(intVal)
				} else {
					field.SetInt(42) // Fallback
				}
			} else {
				field.SetInt(42)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if defaultValue != "" {
				// Convert default string to uint
				if uintVal, err := strconv.ParseUint(defaultValue, 10, 64); err == nil {
					field.SetUint(uintVal)
				} else {
					field.SetUint(42) // Fallback
				}
			} else {
				field.SetUint(42)
			}
		case reflect.Float32, reflect.Float64:
			if defaultValue != "" {
				// Convert default string to float
				if floatVal, err := strconv.ParseFloat(defaultValue, 64); err == nil {
					field.SetFloat(floatVal)
				} else {
					field.SetFloat(42.5) // Fallback
				}
			} else {
				field.SetFloat(42.5)
			}
		case reflect.Bool:
			if defaultValue != "" {
				// Convert default string to bool
				if boolVal, err := strconv.ParseBool(defaultValue); err == nil {
					field.SetBool(boolVal)
				} else {
					field.SetBool(true) // Fallback
				}
			} else {
				field.SetBool(true)
			}
		case reflect.Slice:
			// Create a sample slice with one element
			sliceType := field.Type().Elem()
			sampleSlice := reflect.MakeSlice(field.Type(), 1, 1)

			// If it's a slice of struct, populate the struct with defaults too
			if sliceType.Kind() == reflect.Struct {
				populateSampleValues(sampleSlice.Index(0))
			} else if sliceType.Kind() == reflect.String {
				// For strings, use field name in example
				if defaultValue != "" {
					sampleSlice.Index(0).SetString(defaultValue)
				} else {
					sampleSlice.Index(0).SetString("Sample " + fieldName + " string")
				}
			} else if sliceType.Kind() == reflect.Int {
				if defaultValue != "" {
					// Try to parse default as int
					if intVal, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
						sampleSlice.Index(0).SetInt(intVal)
					} else {
						sampleSlice.Index(0).SetInt(42)
					}
				} else {
					sampleSlice.Index(0).SetInt(42)
				}
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

// DebugLLMInteraction prints detailed debug information about an LLM interaction
// This can be called by any processor when debug is enabled
func DebugLLMInteraction(prompt string, response interface{}) {
	fmt.Println("DEBUG - LLM Prompt:")
	fmt.Println("====================================")
	fmt.Println(prompt)
	fmt.Println("====================================")

	fmt.Printf("DEBUG - Raw LLM Response: %+v\n", response)

	// If response is a string containing a code block, try to clean and print it
	if strResponse, ok := response.(string); ok {
		if strings.HasPrefix(strResponse, "```") {
			fmt.Println("DEBUG - Detected markdown code block, cleaning...")
			// Remove markdown formatting
			strResponse = strings.TrimPrefix(strResponse, "```json")
			strResponse = strings.TrimPrefix(strResponse, "```")
			endIndex := strings.LastIndex(strResponse, "```")
			if endIndex != -1 {
				strResponse = strResponse[:endIndex]
			}
			strResponse = strings.TrimSpace(strResponse)

			// Try to parse and print the cleaned JSON
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(strResponse), &jsonData); err == nil {
				fmt.Printf("DEBUG - Cleaned JSON: %+v\n", jsonData)
			}
		}
	}
}
