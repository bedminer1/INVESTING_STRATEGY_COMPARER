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

	// for demo purposes, remove dates 'in the future'
	var latestDate time.Time
	if len(userData.NetWorthHistory) > 0 {
		latestDate = userData.NetWorthHistory[len(userData.NetWorthHistory)-1].Date
	}

	if !latestDate.IsZero() {
		if err := h.DB.Exec(`
			DELETE FROM portfolio_records 
			WHERE user_id = ? AND date > ?
		`, userData.UserID, latestDate).Error; err != nil {
			fmt.Printf("error deleting newer records: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Failed to delete newer portfolio records",
			})
		}
	}

	// delete dupes
	if err := h.DB.Exec(`
		DELETE FROM portfolio_records 
		WHERE id NOT IN (
			SELECT MIN(id) 
			FROM portfolio_records 
			WHERE user_id = ? 
			GROUP BY date
		) AND user_id = ?
	`, userData.UserID, userData.UserID).Error; err != nil {
		fmt.Printf("error deleting duplicate dates: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete duplicate dates",
		})
	}

	// delete older records (out of top 500)
	if err := h.DB.Exec(`
        DELETE FROM portfolio_records 
        WHERE user_id = ? 
        AND id NOT IN (
            SELECT id FROM (
                SELECT id FROM portfolio_records 
                WHERE user_id = ? 
                ORDER BY date DESC 
                LIMIT 2600
            ) as subquery
        )
    `, userData.UserID, userData.UserID).Error; err != nil {
		fmt.Printf("error deleting older records: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to delete older portfolio records",
		})
	}

	if len(userData.NetWorthHistory) > 0 {
		for i := 0; i < len(userData.NetWorthHistory); i += 100 { // Batch size 100
			end := i + 100
			if end > len(userData.NetWorthHistory) {
				end = len(userData.NetWorthHistory)
			}
			batch := userData.NetWorthHistory[i:end]
			if err := h.DB.CreateInBatches(batch, 100).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error": "Failed to save portfolio records",
				})
			}
		}
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
