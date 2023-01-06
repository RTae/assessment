package expenses

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseByID(echo echo.Context) error {
	var e Expenses
	id := echo.Param("id")
	if id == "" {
		return echo.JSON(
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
			return echo.JSON(
				http.StatusUnprocessableEntity,
				ErrorResponse{Code: http.StatusUnprocessableEntity, Message: "Param id must be integer"},
			)
		}
		if errMatch != nil {
			return echo.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		match, errMatch = regexp.MatchString("no rows in result set", err.Error())
		if match {
			return echo.JSON(
				http.StatusNotFound,
				ErrorResponse{Code: http.StatusNotFound, Message: "Record not found"},
			)
		}
		if errMatch != nil {
			return echo.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		return echo.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}
	return echo.JSON(http.StatusOK, e)
}

func (h *handler) GetExpenses(echo echo.Context) error {
	var expenses []Expenses

	sql := `
	SELECT id, title, amount, note, tags
	FROM expenses
	`
	rows, err := h.db.Query(sql)
	if err != nil {
		return echo.JSON(
			http.StatusInternalServerError,
			ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
		)
	}
	defer rows.Close()

	for rows.Next() {
		var e Expenses
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			return echo.JSON(
				http.StatusInternalServerError,
				ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()},
			)
		}
		expenses = append(expenses, e)
	}
	return echo.JSON(http.StatusOK, expenses)
}
