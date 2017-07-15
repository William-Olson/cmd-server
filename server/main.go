package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! - from: cmd-server")
	})
	e.Logger.Fatal(e.Start(":7447"))

}
