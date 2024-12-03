package db

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Record struct {
	Price float64
	Date  time.Time
}

type Records []Record

func (r *Records) Get(start, end time.Time, fileName string) {
	db, err := gorm.Open(sqlite.Open(fileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	// migrate schema
	if err := db.AutoMigrate(&Record{}); err != nil {
		panic("failed to migrate schema")
	}

	db.Where("date BETWEEN ? AND ?", start, end).Find(&r)
}

func GetCSV(filename string) Records {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	res := make(Records, len(data))

	for i, row := range data[1:] {
		r := Record{}
		parsedDate, err := time.Parse("2006-1-2", row[0])
		if err != nil {
			panic(err)
		}
		r.Date = parsedDate
		r.Price, err = strconv.ParseFloat(row[1], 64)
		if err != nil {
			panic(err)
		}

		res[i] = r
	}

	return res
}
