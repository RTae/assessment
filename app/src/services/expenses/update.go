package expenses

import (
	"net/http"
	"regexp"

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
		match, errMatch := regexp.MatchString("invalid input syntax", err.Error())
		if match {
			return c.JSON(
				http.StatusUnprocessableEntity,
				ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Param id must be integer"},
			)
		}
		if errMatch != nil {
			return c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		match, errMatch = regexp.MatchString("no rows in result set", err.Error())
		if match {
			return c.JSON(
				http.StatusNotFound,
				ErrorResponse{Code: http.StatusNotFound, Message: "Record not found"},
			)
		}
		if errMatch != nil {
			return c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}

	return c.JSON(http.StatusOK, exp)
}
