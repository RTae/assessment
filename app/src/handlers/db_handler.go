package handlers

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RTae/assessment/app/src/settings"
	_ "github.com/lib/pq"
)

func migrateDB(db *sql.DB) {
	var err error

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func InitDB(settings settings.Config) (*sql.DB, func()) {
	var err error
	var db *sql.DB
	db, err = sql.Open("postgres", settings.DatabaseUrl)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	migrateDB(db)
	log.Println("Database Initialized")

	return db, func() { db.Close() }

}

func MockDatabase(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock, func() { db.Close() }
}
