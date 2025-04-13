package processor

import (
	"encoding/json"
	"reflect"
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
