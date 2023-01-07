//go:build unit

package expenses

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RTae/assessment/app/src/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Run("Should create new expense successfully", func(t *testing.T) {
		// Arrange
		e := echo.New()
		body := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		insertMockRow := mock.NewRows([]string{"id"}).AddRow("1")
		mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(insertMockRow)

		h := handler{db}
		c := e.NewContext(req, res)
		expected := "{\"id\":1,\"title\":\"strawberry smoothie\",\"amount\":79,\"note\":\"night market promotion discount 10 bath\",\"tags\":[\"food\",\"beverage\"]}"

		// Act
		err := h.CreateExpense(c)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
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
			e := echo.New()
			body := tt.body
			req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()

			db, mock, close := handlers.MockDatabase(t)
			defer close()

			insertMockRow := mock.NewRows([]string{"id"}).AddRow("1")
			mock.ExpectQuery("INSERT INTO expenses").WillReturnRows(insertMockRow)

			h := handler{db}
			c := e.NewContext(req, res)
			expected := tt.expected

			// Act
			err := h.CreateExpense(c)

			// Assert
			if assert.NoError(t, err) {
				assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
				assert.Regexp(t, expected, strings.TrimSpace(res.Body.String()))
			}

		})
	}

	t.Run("Should return internal error if can not create new expense", func(t *testing.T) {
		// Arrange
		e := echo.New()
		body := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath", 
			"tags": ["food", "beverage"]
		}`
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		db, mock, close := handlers.MockDatabase(t)
		defer close()

		mock.ExpectQuery("INSERT INTO expenses").WillReturnError(sqlmock.ErrCancelled)

		h := handler{db}
		c := e.NewContext(req, res)
		expected := "{\"statusCode\":500,\"message\":\"canceling query due to user request\"}"

		// Act
		err := h.CreateExpense(c)

		// Assert
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, res.Code)
			assert.Equal(t, expected, strings.TrimSpace(res.Body.String()))
		}

	})
}
