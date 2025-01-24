package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()

	db := InitDB("../../price_data.db")
	h := NewHandler(db)

	e.GET("/strategies", h.handleGetStrategies)
	e.GET("/prices", h.handleGetPrices) 
	e.GET("/paper-trading", h.handlePaperTradingStats)

	e.Logger.Fatal(e.Start(":4000"))
}