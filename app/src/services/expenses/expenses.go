package expenses

import "database/sql"

type handler struct {
	db *sql.DB
}

type Expenses struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float32  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type ErrorResponse struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

func CreateHandler(db *sql.DB) *handler {
	return &handler{db}
}
