//go:build unit

package expenses

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RTae/assessment/app/src/handlers"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseHandler(t *testing.T) {
	t.Run("Should get expense by id successfully", func(t *testing.T) {
		// Arrange
		e := echo.New()
		expenseID := "1"
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		getMockRows := mock.NewRows([]string{"ID", "Title", "Amount", "Note", "Tags"}).
			AddRow(
				"1",
				"strawberry smoothie",
				79.00,
				"night market promotion discount 10 bath",
				pq.Array([]string{"food", "beverage"}),
			)

		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
			WithArgs(expenseID).
			WillReturnRows(getMockRows)

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(expenseID)
		expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

		// Act
		err := h.GetExpenseByID(c)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

	t.Run("Should return unprocessable entity error if expense id is empty", func(t *testing.T) {
		// Arrange
		e := echo.New()
		expenseID := ""
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, _, close := handlers.MockDatabase(t)
		defer close()

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(expenseID)
		expected := "{\"statusCode\":422,\"message\":\"Param id is empty\"}"

		// Act
		err := h.GetExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

	t.Run("Should return unprocessable entity error if expense id is not integer", func(t *testing.T) {
		// Arrange
		e := echo.New()
		expenseID := "dwdw"
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
			WithArgs(expenseID).
			WillReturnError(errors.New("invalid input syntax"))

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(expenseID)
		expected := "{\"statusCode\":422,\"message\":\"Param id must be integer\"}"

		// Act
		err := h.GetExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

	t.Run("Should return not found error if the request expense id is not exist", func(t *testing.T) {
		// Arrange
		e := echo.New()
		expenseID := "dwdw"
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
			WithArgs(expenseID).
			WillReturnError(errors.New("no rows in result set"))

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(expenseID)
		expected := "{\"statusCode\":404,\"message\":\"Record not found\"}"

		// Act
		err := h.GetExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

	t.Run("Should return internal error if can not query expense", func(t *testing.T) {
		// Arrange
		e := echo.New()
		expenseID := "dwdw"
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		mock.ExpectQuery("SELECT (.+) FROM expenses WHERE id = ?").
			WithArgs(expenseID).
			WillReturnError(sqlmock.ErrCancelled)

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(expenseID)
		expected := "{\"statusCode\":500,\"message\":\"canceling query due to user request\"}"

		// Act
		err := h.GetExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

}

func TestGetExpensesHandler(t *testing.T) {
	t.Run("Should get expenses successfully", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		getMockRows := sqlmock.NewRows([]string{"ID", "Title", "Amount", "Note", "Tags"}).
			AddRow(
				"1",
				"strawberry smoothie",
				79.00,
				"night market promotion discount 10 bath",
				pq.Array([]string{"food", "beverage"}),
			).
			AddRow(
				"2",
				"Grill pork",
				100.00,
				"night market promotion discount 50 bath",
				pq.Array([]string{"food"}),
			)

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM expenses").
			WillReturnRows(getMockRows)

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense")
		expected := "[{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]},{\"id\":2,\"title\":\"Grill pork\",\"amount\":100,\"note\":\"night market promotion discount 50 bath\",\"tags\":[\"food\"]}]"

		// Act
		err = h.GetExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

	t.Run("Should return internal error if can not query expenses", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		mock.ExpectQuery("SELECT (.+) FROM expenses").
			WillReturnError(sqlmock.ErrCancelled)

		h := handler{db}
		c := e.NewContext(req, res)
		c.SetPath("/expense")
		expected := "{\"statusCode\":500,\"message\":\"canceling query due to user request\"}"

		// Act
		err = h.GetExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})

}
