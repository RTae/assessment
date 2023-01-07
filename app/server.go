package main

import (
	"context"
	"database/sql"
	"fmt"
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

func printBanner() {
	fmt.Println(
		`
		______     __  __     ______   ______     __   __     ______     ______        ______   ______     ______     ______     __  __     __     __   __     ______        ______     __  __     ______     ______   ______     __    __    
		/\  ___\   /\_\_\_\   /\  == \ /\  ___\   /\ "-.\ \   /\  ___\   /\  ___\      /\__  _\ /\  == \   /\  __ \   /\  ___\   /\ \/ /    /\ \   /\ "-.\ \   /\  ___\      /\  ___\   /\ \_\ \   /\  ___\   /\__  _\ /\  ___\   /\ "-./  \   
		\ \  __\   \/_/\_\/_  \ \  _-/ \ \  __\   \ \ \-.  \  \ \___  \  \ \  __\      \/_/\ \/ \ \  __<   \ \  __ \  \ \ \____  \ \  _"-.  \ \ \  \ \ \-.  \  \ \ \__ \     \ \___  \  \ \____ \  \ \___  \  \/_/\ \/ \ \  __\   \ \ \-./\ \  
		 \ \_____\   /\_\/\_\  \ \_\    \ \_____\  \ \_\\"\_\  \/\_____\  \ \_____\       \ \_\  \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \_\ \_\  \ \_\  \ \_\\"\_\  \ \_____\     \/\_____\  \/\_____\  \/\_____\    \ \_\  \ \_____\  \ \_\ \ \_\ 
		  \/_____/   \/_/\/_/   \/_/     \/_____/   \/_/ \/_/   \/_____/   \/_____/        \/_/   \/_/ /_/   \/_/\/_/   \/_____/   \/_/\/_/   \/_/   \/_/ \/_/   \/_____/      \/_____/   \/_____/   \/_____/     \/_/   \/_____/   \/_/  \/_/ 
																																																											   
		KKGo Assessment
		by RTae
	`,
	)
}

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
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "user" || password == "123qweasdzxc" {
			return true, nil
		}
		return false, nil
	}))

}

func main() {
	settings := settings.Setting()
	database, close := handlers.InitDB(settings)
	defer close()

	e := echo.New()
	e.HideBanner = true
	printBanner()

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
