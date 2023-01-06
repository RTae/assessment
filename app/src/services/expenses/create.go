package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) CreateExpenses(c echo.Context) error {

	var exp Expenses
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()},
		)
	}

	sql := `
	INSERT INTO
		expenses (title, amount, note, tags)
	VALUES
		($1, $2, $3, $4) 
	RETURNING id;
	`
	row := h.db.QueryRow(sql, exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	if err := row.Scan(&exp.ID); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}

	return c.JSON(http.StatusCreated, exp)
}
