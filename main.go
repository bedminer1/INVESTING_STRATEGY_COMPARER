package main

import (
	"os"
	"time"

	"github.com/bedminer1/SnP/db"
	"github.com/bedminer1/SnP/strats"
)

func main() {
	start, _ := time.Parse("2006/01/02", "2014/01/01")
	end, _ := time.Parse("2006/01/02", "2025/01/01")

	r := db.Records{}
	r.Get(start, end)

	DCAInvestor := strats.DCA(1000, r)
	VAInvestor := strats.VA(1000, r)
	DynamicVAInvestor := strats.DynamicVA(1000, r)
	MattressInvestor := strats.Mattress(r)

	// Compare
	mostRecentPrice := r[len(r)-1].Price
	strats.CompareStrats(os.Stdout, mostRecentPrice,
		[]strats.Investor{DCAInvestor, VAInvestor, DynamicVAInvestor, MattressInvestor})
}
