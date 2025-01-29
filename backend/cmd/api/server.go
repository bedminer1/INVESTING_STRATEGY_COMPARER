package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	db := InitDB("../../price_data.db")
	h := NewHandler(db)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST},
	}))

	e.GET("/strategies", h.handleGetStrategies)
	e.GET("/prices", h.handleGetPrices)
	e.GET("/paper-trading", h.handlePaperTradingStats)
	e.GET("/fast-paper-trading", h.handleFastPaperTradingStats)
	e.POST("save-metrics", h.handleSaveMetrics)


	e.Logger.Fatal(e.Start(":4000"))
}
