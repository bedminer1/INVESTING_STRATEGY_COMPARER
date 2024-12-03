package main

import (
	"os"
	"time"

	"github.com/bedminer1/SnP/db"
	"github.com/bedminer1/SnP/strats"
	"github.com/labstack/echo/v4"
)

type envelope map[string]interface{}

func handleGet(c echo.Context) error {
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

	priceRecords := db.Records{}
	priceRecords.Get(start, end, "../../price_data.db")

	// Calculate strategy performance using given priceRecords
	DCARecords := strats.WeeklyRecords{
		Strategy: "DCA",
		Records:  strats.DCA(1000, priceRecords),
	}
	VARecords := strats.WeeklyRecords{
		Strategy: "VA",
		Records:  strats.VA(1000, priceRecords),
	}

	DVAcfg := strats.DynamicVAConfig{
		BottomRatio:          4,
		TopRatio:             5.9,
		ReducingMultiplier:   0.9,
		IncreasingMultiplier: 2.31,
	}
	DynamicVARecords := strats.WeeklyRecords{
		Strategy: "DynamicVA",
		Records:  strats.DynamicVA(1000, priceRecords, DVAcfg),
	}
	BuyLowSellHighRecords := strats.WeeklyRecords{
		Strategy: "BuyLowSellHigh",
		Records:  strats.BuyLowSellHigh(priceRecords),
	}
	MattressRecords := strats.WeeklyRecords{
		Strategy: "Mattress",
		Records:  strats.Mattress(priceRecords),
	}

	results := []strats.WeeklyRecords{
		DCARecords,
		VARecords,
		DynamicVARecords,
		BuyLowSellHighRecords,
		MattressRecords,
	}

	strats.CompareStrats(os.Stdout, priceRecords, results)

	return c.JSON(200, envelope{"results": results})
}
