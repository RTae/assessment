//go:build it

package expenses

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	t.Run("Should create new expense successfully", func(t *testing.T) {
		// setup echo server
		e, settings, close := SetupServer(t)
		PingServer()

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
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEqual(t, 0, exp.ID)
		assert.Equal(t, "strawberry smoothie", exp.Title)
		assert.Equal(t, float32(79.00), exp.Amount)
		assert.Equal(t, "night market promotion discount 10 bath", exp.Note)
		assert.Equal(t, []string{"food", "beverage"}, exp.Tags)

		// teardown echo server
		TeardownServer(t, e, close)

	})
}
