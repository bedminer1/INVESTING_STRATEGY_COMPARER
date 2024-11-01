package strats

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bedminer1/SnP/db"
)

// TODO: pick investing strategies, backtest them and compare performance

// Investor represents assets that an investor can own
type Investor struct {
	Strategy string
	Shares   float64
	Reserves float64
}

// NetWorth calculates and returns the NetWorth of the Investor
func (iv Investor) NetWorth(currSnPPrice float64) float64 {
	res := iv.Reserves + (currSnPPrice * iv.Shares)
	return res
}

// SnPValue calculates and returns the value of SnP shares an Investor owns
func (iv Investor) SnPValue(currSnPPrice float64) float64 {
	res := (currSnPPrice * iv.Shares)
	return res
}

// Strat 1: DCA, 1000 per month
func DCA(amount float64, r db.Records) Investor {
	iv := Investor{Strategy: "DCA"}
	investmentMade := false

	for i, record := range r {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		monthlyBuy := amount / record.Price
		iv.Shares += monthlyBuy
		investmentMade = true
	}

	return iv
}

// Strat 2: Value Averaging, target growth of 1000$ per month,
// start two months later and with a 2000$ fund
func VA(targetGrowth float64, r db.Records) Investor {
	iv := Investor{Strategy: "VA", Reserves: 2000}
	investmentMade := false
	monthsCount := 1

	for i, record := range r[2:] {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		iv.Reserves += 1000
		// calculate delta from target
		snpv := iv.SnPValue(record.Price)
		targetSnPValue := float64(monthsCount) * targetGrowth
		amountToMove := targetSnPValue - snpv

		// check if enough
		if amountToMove <= iv.Reserves {
			iv.Reserves -= amountToMove // also works for adding to it
			iv.Shares += amountToMove / record.Price
		} else {
			iv.Shares += iv.Reserves / record.Price
			iv.Reserves = 0
		}

		investmentMade = true
		monthsCount++
	}

	return iv
}

// Strat 3: VA where targetGrowth scales with reserves
func DynamicVA(targetGrowth float64, r db.Records) Investor {
	iv := Investor{Strategy: "DynamicVA", Reserves: 2000}
	investmentMade := false
	monthsCount := 1
	targetSnPValue := float64(0)

	for i, record := range r[2:] {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		iv.Reserves += 1000
		// calculate delta from target
		snpv := iv.SnPValue(record.Price)
		targetSnPValue += targetGrowth
		amountToMove := targetSnPValue - snpv

		// check if enough
		if amountToMove <= iv.Reserves {
			iv.Reserves -= amountToMove // also works for adding to it
			iv.Shares += amountToMove / record.Price
		} else {
			iv.Shares += iv.Reserves / record.Price
			iv.Reserves = 0
		}

		// Adjust target growth based on cash buffer ratio
		cashBufferRatio := iv.Reserves / targetGrowth
		if cashBufferRatio < 4 {
			targetGrowth = max(targetGrowth * 0.9, 1000) // reduce growth
		} else if cashBufferRatio > 5.9 {
			targetGrowth *= 2.3 // increase growth
		}


		investmentMade = true
		monthsCount++
	}

	return iv
}

// Strat 4: Mattress Stuffer
func Mattress(r db.Records) Investor {
	iv := Investor{Strategy: "Mattress"}
	investmentMade := false

	for i, record := range r {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		iv.Reserves += 1000
		investmentMade = true
	}

	return iv
}

func CompareStrats(out io.Writer, currSnPPrice float64, investors []Investor) {
	fmt.Fprintln(out, "Number of Investors:", len(investors))

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Strategy\tNet Worth\tSnP Value\tReverses Value\t")

	for _, investor := range investors {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t\n", 
		investor.Strategy, 
		investor.NetWorth(currSnPPrice),
		investor.SnPValue(currSnPPrice),
		investor.Reserves,
		)
	}
	fmt.Fprintln(w)

	w.Flush()
}
