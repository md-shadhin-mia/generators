package generators

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// GenerateJSONSchema generates a JSON Schema validator for a GORM model using reflection.
func GenerateJSONSchema(model interface{}) (string, error) {
	type schemaProperty struct {
		Type    string `json:"type,omitempty"`
		Format  string `json:"format,omitempty"`
		Minimum *int   `json:"minimum,omitempty"`
	}

	type jsonSchema struct {
		Type       string                    `json:"type"`
		Properties map[string]schemaProperty `json:"properties"`
		Required   []string                  `json:"required,omitempty"`
	}

	schema := jsonSchema{
		Type:       "object",
		Properties: make(map[string]schemaProperty),
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("model must be a struct or a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "" {
			continue
		}

		fieldType := field.Type.Kind()
		property := schemaProperty{}

		switch fieldType {
		case reflect.String:
			property.Type = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			property.Type = "integer"
			if minValue, ok := field.Tag.Lookup("min"); ok {
				min, err := strconv.Atoi(minValue)
				if err == nil {
					property.Minimum = &min
				}
			}
		case reflect.Float32, reflect.Float64:
			property.Type = "number"
		case reflect.Bool:
			property.Type = "boolean"
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				property.Type = "array"
				property.Format = "string"
			}
		}

		schema.Properties[jsonTag] = property

		if binding, ok := field.Tag.Lookup("binding"); ok && binding == "required" {
			schema.Required = append(schema.Required, jsonTag)
		}
	}

	output, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}
