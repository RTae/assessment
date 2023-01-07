//go:build it

package expenses

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	// setup echo server
	e, settings, close := SetupServer(t)
	PingServer()

	t.Run("Should create new expense successfully", func(t *testing.T) {
		// Arrange
		body := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`
		var exp Expenses

		// Act
		res := Request(t, http.MethodPost, Uri(fmt.Sprint(settings.Port), "expenses"), strings.NewReader(body))
		err := res.Decode(&exp)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, res.StatusCode)
			assert.NotEqual(t, 0, exp.ID)
			assert.Equal(t, "strawberry smoothie", exp.Title)
			assert.Equal(t, float32(79.00), exp.Amount)
			assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
			assert.Equal(t, []string{"food", "beverage"}, exp.Tags)
		}

	})

	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			"Should return unprocess entity error if title is not correct",
			`{
				"title": 213,
				"amount": 79,
				"note": "night market promotion discount 10 bath", 
				"tags": ["food", "beverage"]
			}`,
			"cannot unmarshal number into Go struct field Expenses.title of type string",
		},
		{
			"Should return unprocess entity error if amount is not correct",
			`{
				"title": "strawberry smoothie",
				"amount": "79",
				"note": "night market promotion discount 10 bath", 
				"tags": ["food", "beverage"]
			}`,
			"cannot unmarshal string into Go struct field Expenses.amount of type float32",
		},
		{
			"Should return unprocess entity error if note is not correct",
			`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": 22321, 
				"tags": ["food", "beverage"]
			}`,
			"cannot unmarshal number into Go struct field Expenses.note of type string",
		},
		{
			"Should return unprocess entity error if tags is not correct",
			`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath", 
				"tags": "["food", "beverage"]"
			}`,
			"invalid character 'f' after object key:value pair",
		},
		{
			"Should return unprocess entity error if data in tags is not correct",
			`{
				"title": "strawberry smoothie",
				"amount": 79,
				"note": "night market promotion discount 10 bath", 
				"tags": ["food", "beverage", 2312]
			}`,
			"cannot unmarshal number into Go struct field Expenses.tags of type string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			body := tt.body
			expected := tt.expected
			var errRes ErrorResponse

			// Act
			res := Request(t, http.MethodPost, Uri(fmt.Sprint(settings.Port), "expenses"), strings.NewReader(body))
			err := res.Decode(&errRes)

			// Assert
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusUnprocessableEntity, errRes.Code)
				assert.Regexp(t, expected, errRes.Message)
			}

		})
	}

	// teardown echo server
	TeardownServer(t, e, close)
}
