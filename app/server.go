package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RTae/assessment/app/src/handlers"
	"github.com/RTae/assessment/app/src/services/expenses"
	"github.com/RTae/assessment/app/src/settings"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func initRoute(e *echo.Echo, db *sql.DB) {

	expensesHandler := expenses.CreateHandler(db)

	g := e.Group("expenses")
	g.POST("", expensesHandler.CreateExpense)
	g.GET("/:id", expensesHandler.GetExpenseByID)
	g.PUT("/:id", expensesHandler.UpdateExpenseByID)
	g.GET("", expensesHandler.GetExpenses)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
}

func initMiddleware(e *echo.Echo, db *sql.DB) {
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

}

func main() {
	settings := settings.Setting()
	database, close := handlers.InitDB(settings)
	defer close()

	e := echo.New()

	initMiddleware(e, database)
	initRoute(e, database)

	go func() {
		if err := e.Start(settings.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Print("Server stopped")
}
