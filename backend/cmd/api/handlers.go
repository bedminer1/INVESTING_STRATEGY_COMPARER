package main

import (
	"time"

	"github.com/bedminer1/SnP/models"
	"github.com/bedminer1/SnP/strats"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type envelope map[string]interface{}

type Handler struct {
	DB *gorm.DB
}

func InitDB(fileName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	// migrate schema
	if err := db.AutoMigrate(&models.Record{}); err != nil {
		panic("failed to migrate schema")
	}

	return db
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) fetchPrices(start, end time.Time) *models.Records {
	priceRecords := &models.Records{}
	h.DB.Where("date BETWEEN ? AND ?", start, end).Find(&priceRecords)

	return priceRecords
}

func (h *Handler) handleGetPrices(c echo.Context) error {
	startInput := c.QueryParam("start")
	endInput := c.QueryParam("end")
	if startInput == "" {
		startInput = "2014_01_01"
	}
	if endInput == "" {
		endInput = "2025_01_01"
	}
	start, _ := time.Parse("2006_01_02", startInput)
	end, _ := time.Parse("2006_01_02", endInput)

	priceRecords := *h.fetchPrices(start, end)
	return c.JSON(200, envelope{"price_data": priceRecords})
}

func (h *Handler) handleGetStrategies(c echo.Context) error {
	startInput := c.QueryParam("start")
	endInput := c.QueryParam("end")
	if startInput == "" {
		startInput = "2014_01_01"
	}
	if endInput == "" {
		endInput = "2025_01_01"
	}
	start, _ := time.Parse("2006_01_02", startInput)
	end, _ := time.Parse("2006_01_02", endInput)

	priceRecords := *h.fetchPrices(start, end)

	// Calculate strategy performance using given priceRecords
	DCARecords := models.WeeklyRecords{
		Strategy: "DCA",
		Records:  strats.DCA(1000, priceRecords),
	}
	VARecords := models.WeeklyRecords{
		Strategy: "VA",
		Records:  strats.VA(1000, priceRecords),
	}

	DVAcfg := strats.DynamicVAConfig{
		BottomRatio:          4,
		TopRatio:             5.9,
		ReducingMultiplier:   0.9,
		IncreasingMultiplier: 2.31,
	}
	DynamicVARecords := models.WeeklyRecords{
		Strategy: "DynamicVA",
		Records:  strats.DynamicVA(1000, priceRecords, DVAcfg),
	}
	BuyLowSellHighRecords := models.WeeklyRecords{
		Strategy: "BuyLowSellHigh",
		Records:  strats.BuyLowSellHigh(priceRecords),
	}
	MattressRecords := models.WeeklyRecords{
		Strategy: "Mattress",
		Records:  strats.Mattress(priceRecords),
	}

	results := []models.WeeklyRecords{
		DCARecords,
		VARecords,
		DynamicVARecords,
		BuyLowSellHighRecords,
		MattressRecords,
	}

	return c.JSON(200, envelope{"results": results})
}
