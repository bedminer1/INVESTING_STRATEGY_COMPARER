package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()

	e.GET("/", handleGet)

	e.Logger.Fatal(e.Start(":4000"))
}