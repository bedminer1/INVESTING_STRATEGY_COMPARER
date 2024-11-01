package strats

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"

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

	for i, record := range r {
		if i < 2 {
			continue
		}
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

type DynamicVAConfig struct {
	BottomRatio          float64
	TopRatio             float64
	ReducingMultiplier   float64
	IncreasingMultiplier float64
}

// Strat 3: VA where targetGrowth scales with reserves
func DynamicVA(targetGrowth float64, r db.Records, cfg DynamicVAConfig) Investor {
	iv := Investor{Strategy: "DynamicVA", Reserves: 2000}
	investmentMade := false
	monthsCount := 1
	targetSnPValue := float64(0)

	for i, record := range r {
		if i < 2 {
			continue
		}
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
		if cashBufferRatio < cfg.BottomRatio {
			targetGrowth = max(targetGrowth*cfg.ReducingMultiplier, 1000) // reduce growth
		} else if cashBufferRatio > cfg.TopRatio {
			targetGrowth *= cfg.IncreasingMultiplier // increase growth
		}

		investmentMade = true
		monthsCount++
	}

	return iv
}

// calculateAverage computes the average of a slice of float64 numbers
func calculateAverage(prices []float64) float64 {
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	return sum / float64(len(prices))
}

// Strat 4: Buy Low Sell High
// Noted that this strategy has potential for good performance but requires very precise tuning, like DVA
// Small tweaks in the buy, sell threshold and windowsize lead to big differences in performance
func BuyLowSellHigh(r db.Records) Investor {
	iv := Investor{Strategy: "BuyLowSellHigh", Reserves: 0} // Starting with a reserve balance
	windowSize := 20
	buyThreshold := 1.00
	sellThreshold := 1.03
	investmentMade := false

	var recentPrices []float64

	for i, record := range r {
		// Add the current price to recent prices
		recentPrices = append(recentPrices, record.Price)
		if len(recentPrices) > windowSize {
			recentPrices = recentPrices[1:] // Keep only the last `windowSize` prices
		}

		// Calculate the moving average of recent prices
		if len(recentPrices) < windowSize {
			continue // Not enough data points for the window yet
		}

		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		investmentMade = true

		avgPrice := calculateAverage(recentPrices)
		iv.Reserves += 1000

		// Buy condition
		if record.Price < buyThreshold*avgPrice {
			sharesToBuy := iv.Reserves / record.Price
			iv.Shares += sharesToBuy
			iv.Reserves -= sharesToBuy * record.Price
		}

		// Sell
		if record.Price > sellThreshold*avgPrice {
			proceeds := iv.Shares * record.Price
			iv.Reserves += proceeds
			iv.Shares = 0
		}
	}

	return iv
}

// Strat 5: Mattress Stuffer
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

func CompareStrats(out io.Writer, currSnPPrice float64, investors []Investor, start, end time.Time) {
	fmt.Fprintf(out, "Backtested from %s to %s\n", start.Format("2 Jan 2006"), end.Format("2 Jan 2006"))
	fmt.Fprintln(out, "Number of Investors:", len(investors))
	fmt.Fprintln(out)


	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Strategy\tNet Worth\tSnP Value\tNo of Shares\tReverses Value\t")

	for _, investor := range investors {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t\n",
			investor.Strategy,
			investor.NetWorth(currSnPPrice),
			investor.SnPValue(currSnPPrice),
			investor.Shares,
			investor.Reserves,
		)
	}
	fmt.Fprintln(w)

	w.Flush()
}
