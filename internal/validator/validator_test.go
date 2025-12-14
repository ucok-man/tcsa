package validator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_CustomPortTranslation(t *testing.T) {
	v := New()

	type Config struct {
		Port uint `validate:"port"`
	}

	t.Run("should apply custom port translation format", func(t *testing.T) {
		err := v.Struct(Config{Port: 99999})

		require.Error(t, err)
		errMsg := err.Error()

		assert.Contains(t, errMsg, "has invalid value of")
		assert.Contains(t, errMsg, "99999")
	})
}

func TestValidationErrorMap_Error(t *testing.T) {
	tests := []struct {
		name     string
		errMap   ValidationErrorMap
		expected string
	}{
		{
			name:     "empty map",
			errMap:   ValidationErrorMap{},
			expected: "map[]",
		},
		{
			name:     "single error",
			errMap:   ValidationErrorMap{"Email": "invalid email format"},
			expected: "Email: invalid email format",
		},
		{
			name: "multiple errors contain all fields",
			errMap: ValidationErrorMap{
				"Email": "invalid email",
				"Name":  "required field",
			},
			expected: "Email: invalid email; Name: required field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errMap.Error()

			if tt.expected != "" {
				assert.Equal(t, tt.expected, result)
			}

			if len(tt.errMap) > 1 {
				assert.Contains(t, result, ";")

				results := strings.Split(result, ";")
				assert.Equal(t, len(results), 2)

			}
		})
	}
}

func TestValidationErrorMap_MarshalJSON(t *testing.T) {
	t.Run("strips struct name and lowercases keys", func(t *testing.T) {
		errMap := ValidationErrorMap{
			"User.Email": "invalid email format",
			"User.Name":  "required field",
		}

		data, err := json.Marshal(errMap)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, "invalid email format", result["email"])
		assert.Equal(t, "required field", result["name"])
	})

	t.Run("handles nested struct paths", func(t *testing.T) {
		errMap := ValidationErrorMap{
			"User.Address.City": "required field",
		}

		data, err := json.Marshal(errMap)
		require.NoError(t, err)

		var result map[string]string
		err = json.Unmarshal(data, &result)
		require.NoError(t, err)

		assert.Equal(t, "required field", result["address.city"])
	})

	t.Run("handles empty map", func(t *testing.T) {
		errMap := ValidationErrorMap{}

		data, err := json.Marshal(errMap)
		require.NoError(t, err)

		assert.Equal(t, "{}", string(data))
	})
}
