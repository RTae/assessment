package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) UpdateExpenseByID(c echo.Context) error {
	exp := new(Expenses)
	expenseId := c.Param("id")
	if expenseId == "" {
		return c.JSON(
			http.StatusUnprocessableEntity,
			ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Param id is empty"},
		)
	}

	if err := c.Bind(exp); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Invalid request body"},
		)
	}

	sql := `
	UPDATE 
		expenses SET title = $1, amount = $2, note = $3, tags = $4
	WHERE
		id = $5
	RETURNING id
	`
	row := h.db.QueryRow(sql, exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags), expenseId)

	err := row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, exp)
}
