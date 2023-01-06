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
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpense(t *testing.T) {
	t.Run("Should update expense successfully", func(t *testing.T) {
		// Arrange
		updateExpenseID := "3"
		body := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount", 
			"tags": ["beverage"]
		}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expense", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		resultMockRow := mock.NewRows([]string{"ID"}).AddRow(updateExpenseID)
		mock.ExpectQuery("UPDATE expenses").
			WillReturnRows(resultMockRow)

		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(updateExpenseID)
		expected := "{\"id\":3,\"title\":\"apple smoothie\",\"amount\":89,\"note\":\"no discount\",\"tags\":[\"beverage\"]}"

		// Act
		err := h.UpdateExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})

	t.Run("Should return unprocessable entity error if expense id is empty", func(t *testing.T) {
		// Arrange
		updateExpenseID := ""
		body := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount", 
			"tags": ["beverage"]
		}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expense", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, _, close := handlers.MockDatabase(t)
		defer close()

		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(updateExpenseID)
		expected := "{\"statusCode\":422,\"message\":\"Param id is empty\"}"

		// Act
		err := h.UpdateExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})

	t.Run("Should return unprocessable entity error if expense id is not integer", func(t *testing.T) {
		// Arrange
		e := echo.New()
		updateExpenseID := "dw2"
		body := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount", 
			"tags": ["beverage"]
		}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		resultMockRow := mock.NewRows([]string{"ID"}).AddRow(updateExpenseID)
		mock.ExpectQuery("UPDATE expenses").
			WillReturnRows(resultMockRow).
			WillReturnError(errors.New("invalid input syntax"))

		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(updateExpenseID)
		expected := "{\"statusCode\":422,\"message\":\"Param id must be integer\"}"

		// Act
		err := h.UpdateExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})

	t.Run("Should return not found error if the request expense id is not exist", func(t *testing.T) {
		// Arrange
		e := echo.New()
		updateExpenseID := "9"
		body := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount", 
			"tags": ["beverage"]
		}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		resultMockRow := mock.NewRows([]string{"ID"}).AddRow(updateExpenseID)
		mock.ExpectQuery("UPDATE expenses").
			WillReturnRows(resultMockRow).
			WillReturnError(errors.New("no rows in result set"))

		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(updateExpenseID)
		expected := "{\"statusCode\":404,\"message\":\"Record not found\"}"

		// Act
		err := h.UpdateExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})

	t.Run("Should return internal error if can not query expense", func(t *testing.T) {
		// Arrange
		e := echo.New()
		updateExpenseID := "3"
		body := `{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount", 
			"tags": ["beverage"]
		}`
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		resultMockRow := mock.NewRows([]string{"ID"}).AddRow(updateExpenseID)
		mock.ExpectQuery("UPDATE expenses").
			WillReturnRows(resultMockRow).
			WillReturnError(sqlmock.ErrCancelled)

		h := handler{db}
		c := e.NewContext(req, rec)
		c.SetPath("/expense/:id")
		c.SetParamNames("id")
		c.SetParamValues(updateExpenseID)
		expected := "{\"statusCode\":500,\"message\":\"canceling query due to user request\"}"

		// Act
		err := h.UpdateExpenseByID(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}

	})
}
