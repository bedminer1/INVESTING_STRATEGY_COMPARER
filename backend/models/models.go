package models

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	Price float64
	Date  time.Time
}

type PortfolioRecord struct {
	gorm.Model
	UserID string    `json:"user_id" gorm:"index"`
	Price  float64   `json:"price"`
	Date   time.Time `json:"date"`
}

type WeeklyRecord struct {
	Time     time.Time
	NetWorth float64
	SnpValue float64
	Shares   float64
	Reserves float64
}

type WeeklyRecords struct {
	Strategy string
	Records  []WeeklyRecord
}

type User struct {
	UserID          string            `json:"user_id" gorm:"uniqueIndex"`
	Cash            float64           `json:"cash"`
	Position        float64           `json:"position"`
	NetWorthHistory []PortfolioRecord `json:"net_worth_history" gorm:"foreignKey:UserID;references:UserID"`
}
