//go:build it

package expenses

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExpenseByID(t *testing.T) {
	// setup echo server
	e, settings, close := SetupServer(t)
	PingServer()
	t.Run("Should get expense successfully", func(t *testing.T) {

		// Arrange
		createExpense := SeedExpense(t, settings)

		// Act
		var exp Expenses
		res := Request(t, http.MethodGet, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%d", createExpense.ID)), nil)
		err := res.Decode(&exp)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, createExpense.ID, exp.ID)
			assert.Equal(t, createExpense.Title, exp.Title)
			assert.Equal(t, createExpense.Amount, exp.Amount)
			assert.Equal(t, createExpense.Note, exp.Note)
			assert.Equal(t, createExpense.Tags, exp.Tags)
		}
	})

	t.Run("Should return unprocessable entity error if expense id is not integer", func(t *testing.T) {

		// Arrange
		expenseId := "12d"
		expectedMessage := "Param id must be integer"

		// Act
		var errRes ErrorResponse
		res := Request(t, http.MethodGet, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%s", expenseId)), nil)
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
		expectedMessage := "Record not found"

		// Act
		var errRes ErrorResponse
		res := Request(t, http.MethodGet, Uri(fmt.Sprint(settings.Port), fmt.Sprintf("expenses/%s", expenseId)), nil)
		err := res.Decode(&errRes)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, errRes.Code)
			assert.Equal(t, expectedMessage, errRes.Message)
		}
	})

	// teardown echo server
	TeardownServer(t, e, close)
}

func TestGetExpenses(t *testing.T) {
	// setup echo server
	e, settings, close := SetupServer(t)
	PingServer()
	t.Run("Should get expenses successfully", func(t *testing.T) {

		// Arrange
		SeedExpense(t, settings)

		// Act
		var exp []Expenses
		res := Request(t, http.MethodGet, Uri(fmt.Sprint(settings.Port), "expenses"), nil)
		err := res.Decode(&exp)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.NotEqual(t, 0, len(exp))
		}
	})

	// teardown echo server
	TeardownServer(t, e, close)
}
