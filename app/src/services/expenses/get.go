package expenses

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseByID(c echo.Context) error {
	var e Expenses
	id := c.Param("id")
	if id == "" {
		return c.JSON(
			http.StatusUnprocessableEntity,
			ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Param id is empty"},
		)
	}

	sql := `
	SELECT id, title, amount, note, tags
	FROM expenses
	WHERE id = $1
	`
	err := h.db.QueryRow(sql, id).Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

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
	return c.JSON(http.StatusOK, e)
}

func (h *handler) GetExpenses(c echo.Context) error {
	var expenses []Expenses

	sql := `
	SELECT id, title, amount, note, tags
	FROM expenses
	`
	rows, err := h.db.Query(sql)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}
	defer rows.Close()

	for rows.Next() {
		var e Expenses
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		expenses = append(expenses, e)
	}
	return c.JSON(http.StatusOK, expenses)
}
