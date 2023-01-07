//go:build it

package expenses

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseByID(t *testing.T) {
	// setup echo server
	e, settings, close := SetupServer(t)
	PingServer()
	t.Run("Should update expense successfully", func(t *testing.T) {

		// Arrange
		createExpense := SeedExpense(t, settings)

		// Act
		body := `{
			"title":"apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`

		var exp Expenses
		res := Request(t, http.MethodPut, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%d", createExpense.ID)), strings.NewReader(body))
		err := res.Decode(&exp)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, exp.ID, createExpense.ID)
			assert.Equal(t, exp.Title, "apple smoothie")
			assert.Equal(t, exp.Amount, float32(89))
			assert.Equal(t, exp.Note, "no discount")
			assert.Equal(t, exp.Tags, []string{"beverage"})
		}
	})

	t.Run("Should return unprocessable entity error if expense id is not integer", func(t *testing.T) {

		// Arrange
		expenseId := "12d"
		body := `{
			"title":"apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`
		expectedMessage := "Param id must be integer"

		// Act
		var errRes ErrorResponse
		res := Request(t, http.MethodPut, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%s", expenseId)), strings.NewReader(body))
		err := res.Decode(&errRes)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, errRes.Code)
			assert.Equal(t, expectedMessage, errRes.Message)
		}
	})

	t.Run("Should return not found error if the request expense id is not exist", func(t *testing.T) {

		// Arrange
		expenseId := "1000"
		body := `{
			"title":"apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`
		expectedMessage := "Record not found"

		// Act
		var errRes ErrorResponse
		res := Request(t, http.MethodPut, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%s", expenseId)), strings.NewReader(body))
		err := res.Decode(&errRes)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, errRes.Code)
			assert.Equal(t, expectedMessage, errRes.Message)
		}
	})

	t.Run("Should return unprocess entity error if body is not correct", func(t *testing.T) {

		// Arrange
		expenseId := "1000"
		body := `{
			"title":12323,
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`
		expectedMessage := "Invalid request body"

		// Act
		var errRes ErrorResponse
		res := Request(t, http.MethodPut, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%s", expenseId)), strings.NewReader(body))
		err := res.Decode(&errRes)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, errRes.Code)
			assert.Equal(t, expectedMessage, errRes.Message)
		}
	})

	// teardown echo server
	TeardownServer(t, e, close)
}
