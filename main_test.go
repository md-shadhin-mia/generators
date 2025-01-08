package generators_test

import (
	"md-shadhin-mia/generators"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJSONSchema(t *testing.T) {
	t.Run("Generate JSON schema for a simple struct", func(t *testing.T) {
		// Add test code here
		type User struct {
			Name     string   `json:"name" binding:"required"`
			Age      int      `json:"age" binding:"required" min:"0"`
			Email    string   `json:"email"`
			IsActive bool     `json:"is_active"`
			Tags     []string `json:"tags"`
		}

		schema, err := generators.GenerateJSONSchema(User{})
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		expected := `{
  "type": "object",
  "properties": {
    "age": {
      "type": "integer",
      "minimum": 0
    },
    "email": {
      "type": "string"
    },
    "is_active": {
      "type": "boolean"
    },
    "name": {
      "type": "string"
    },
    "tags": {
      "type": "array",
      "format": "string"
    }
  },
  "required": [
    "name",
    "age"
  ]
}`
		if schema != expected {
			t.Fatalf("unexpected schema:\n%s", schema)
		}

	})

	t.Run("Basic Struct with String and Integer Fields", func(t *testing.T) {
		type TestModel struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		expected := `{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "age": {
      "type": "integer"
    }
  }
}`

		result, err := generators.GenerateJSONSchema(TestModel{})
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Struct with Required Fields", func(t *testing.T) {
		type TestModel struct {
			Name string `json:"name" binding:"required"`
			Age  int    `json:"age"`
		}

		expected := `{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "age": {
      "type": "integer"
    }
  },
  "required": ["name"]
}`

		result, err := generators.GenerateJSONSchema(TestModel{})
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Struct with Minimum Value Constraint", func(t *testing.T) {
		type TestModel struct {
			Age int `json:"age" min:"18"`
		}

		expected := `{
  "type": "object",
  "properties": {
    "age": {
      "type": "integer",
      "minimum": 18
    }
  }
}`

		result, err := generators.GenerateJSONSchema(TestModel{})
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Invalid Input (Non-Struct)", func(t *testing.T) {
		invalidInput := 42

		_, err := generators.GenerateJSONSchema(invalidInput)
		assert.Error(t, err)
		assert.Equal(t, "model must be a struct or a pointer to a struct", err.Error())
	})

	t.Run("Struct with Ignored Field", func(t *testing.T) {
		type TestModel struct {
			Name string `json:"-"`
			Age  int    `json:"age"`
		}

		expected := `{
  "type": "object",
  "properties": {
    "age": {
      "type": "integer"
    }
  }
}`

		result, err := generators.GenerateJSONSchema(TestModel{})
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

}
