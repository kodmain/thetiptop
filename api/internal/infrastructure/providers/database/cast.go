package database

import "reflect"

// StructToMap Convert a struct to a map[string]any without using marshal
// The function uses reflection to analyze fields and their tags.
//
// Parameters:
// - input: any The input struct.
//
// Returns:
// - result: map[string]any The resulting map with field names as keys and values as values.
func StructToMap(input any) map[string]any {
	result := make(map[string]any)

	// Get the reflection object of the struct
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)

	// Handle pointer to struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Use the JSON tag as the key if it exists
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = field.Name
		}

		// Check if the field value is a pointer and handle nil pointers
		if value.Kind() == reflect.Ptr && value.IsNil() {
			result[tag] = nil
		} else if value.Kind() == reflect.Ptr {
			result[tag] = value.Elem().Interface()
		} else {
			result[tag] = value.Interface()
		}
	}

	return result
}
