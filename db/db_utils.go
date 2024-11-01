package db

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Record struct {
	Price float64
	Date  time.Time
}

type Records []Record

func (r *Records) Get(start, end time.Time) {
	db, err := gorm.Open(sqlite.Open("price_data.db"), &gorm.Config{})
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

func (r Records) List(out io.Writer, limit int) {
	fmt.Fprintln(out, "Number of Records Fetched:", len(r))

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Date\tPrice\t")

	if limit > len(r)-1 {
		limit = len(r)
	}

	if limit == 0 {
		fmt.Fprintln(out)
		return
	}

	for _, record := range r[:limit] {
		fmt.Fprintf(w, "%s\t%.4f\t\n", record.Date.Format("2006/01/02"), record.Price)
	}
	fmt.Fprintln(w)

	w.Flush()
}
