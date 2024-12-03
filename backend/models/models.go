package models

import "time"

type Record struct {
	Price float64
	Date  time.Time
}

type Records []Record

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