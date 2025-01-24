package main

import (
	"fmt"
	"time"

	"github.com/bedminer1/SnP/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) handlePaperTradingStats(c echo.Context) error {
	records := []models.Record{}
	h.DB.Find(&records)

	numRecords := len(records)

	now := time.Now()
	currentTimeOfDay := now.Hour()*3600 + now.Minute()*60 + now.Second() // Time of day in seconds
	totalSeconds := int((time.Hour * 24).Seconds())                        // Total simulation seconds
	progress := float64(currentTimeOfDay) / float64(totalSeconds)
	currentIndex := int(progress * float64(numRecords))

	if currentIndex > numRecords {
		currentIndex = numRecords
	}

	return c.JSON(200, echo.Map{
		"records": records[:currentIndex],
		"progress": fmt.Sprintf("%d / %d", currentIndex, numRecords),
	})
}
