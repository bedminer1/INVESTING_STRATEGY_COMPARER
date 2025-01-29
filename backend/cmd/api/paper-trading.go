package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bedminer1/SnP/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) handlePaperTradingStats(c echo.Context) error {
	records := []models.Record{}
	h.DB.Find(&records)

	numRecords := len(records)

	now := time.Now()
	currentTimeOfDay := now.Hour()*3600 + now.Minute()*60 + now.Second() // Time of day in seconds
	totalSeconds := int((time.Hour * 24).Seconds())                      // Total simulation seconds
	progress := float64(currentTimeOfDay) / float64(totalSeconds)
	currentIndex := int(progress * float64(numRecords))

	if currentIndex > numRecords {
		currentIndex = numRecords
	}

	return c.JSON(200, echo.Map{
		"records":  records[:currentIndex],
		"progress": fmt.Sprintf("%d / %d", currentIndex, numRecords),
	})
}

func (h *Handler) handleFastPaperTradingStats(c echo.Context) error {
	records := []models.Record{}
	h.DB.Find(&records)

	numRecords := len(records)

	now := time.Now()
	minutes := now.Minute()
	seconds := now.Second()
	totalSeconds := 60 * 60
	elapsedSeconds := (minutes * 60) + seconds

	progress := float64(elapsedSeconds) / float64(totalSeconds)
	currentIndex := int(progress * float64(numRecords))

	if currentIndex > numRecords {
		currentIndex = numRecords
	}

	return c.JSON(200, echo.Map{
		"records":  records[:currentIndex],
		"progress": fmt.Sprintf("%d / %d", currentIndex, numRecords),
	})
}

func (h *Handler) handleSaveUserStats(c echo.Context) error {
	var userData models.User
	if err := c.Bind(&userData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request data: " + err.Error(),
		})
	}

	var existingUser models.User
	result := h.DB.Where("user_id = ?", userData.UserID).First(&existingUser)

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		h.DB.Create(&userData)
	} else {
		h.DB.Model(&existingUser).Where("user_id = ?", userData.UserID).Updates(userData)
	}

	return c.JSON(200, echo.Map{
		"message":   "user stats saved",
		"user_info": userData,
	})
}

func (h *Handler) handleFetchUserStats(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "user_id is required",
		})
	}

	var user models.User
	if err := h.DB.Preload("NetWorthHistory").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user_info": user,
	})
}
