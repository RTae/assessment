package expenses

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/RTae/assessment/app/src/handlers"
	"github.com/RTae/assessment/app/src/settings"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func SetupServer(t *testing.T) (*echo.Echo, settings.Config, func()) {
	e := echo.New()
	var settings = settings.Setting()
	database, close := handlers.InitDB(settings)

	go func(c *echo.Echo) {
		expensesHandler := CreateHandler(database)

		g := c.Group("expenses")
		g.POST("", expensesHandler.CreateExpense)
		g.GET("/:id", expensesHandler.GetExpenseByID)
		g.PUT("/:id", expensesHandler.UpdateExpenseByID)
		g.GET("", expensesHandler.GetExpenses)

		c.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "OK")
		})

		c.Start(fmt.Sprintf("%s%s", settings.Url, settings.Port))
	}(e)

	return e, settings, close
}

func PingServer() {
	for {
		var settings = settings.Setting()
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s%s", settings.Url, settings.Port), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func TeardownServer(t *testing.T, eh *echo.Echo, close func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer close()
	err := eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func (res *Response) Decode(v interface{}) error {
	if res.err != nil {
		return res.err
	}
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(result, v)
}

func Request(t *testing.T, method, url string, body io.Reader) *Response {

	if body == nil {
		body = bytes.NewBufferString("")
	}

	req, err := http.NewRequest(method, url, body)
	assert.NoError(t, err)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := http.Client{}
	resp, err := client.Do(req)

	assert.NoError(t, err)

	return &Response{resp, err}
}

func Uri(port string, paths ...string) string {
	host := "http://localhost" + port
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func SeedExpense(t *testing.T, settings settings.Config) Expenses {
	body := `{
		"title":"pay market",
		"amount": 9999.00,
		"note": "clear debt",
		"tags": ["markets", "debt"]
	}`

	var createExpense Expenses
	res := Request(t, http.MethodPost, Uri(fmt.Sprint(settings.Port), "expenses"), strings.NewReader(body))
	err := res.Decode(&createExpense) // ใช้ decode ข้อมูลที่ได้จาก response body มาเก็บไว้ในตัวแปร u
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}

	return createExpense
}
