package strats

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/bedminer1/SnP/models"
)

// NetWorth calculates and returns the NetWorth of the Investor
func CalculateNetWorth(currSnPPrice, shares, reserves float64) float64 {
	res := reserves + (currSnPPrice * shares)
	return res
}

// SnPValue calculates and returns the value of SnP shares an Investor owns
func CalculateSnPValue(currSnPPrice, shares, reserves float64) float64 {
	res := (currSnPPrice * shares)
	return res
}

// Strat 1: DCA, 1000 per month
func DCA(amount float64, priceRecords models.Records) []models.WeeklyRecord {
	investmentMade := false
	records := []models.WeeklyRecord{}
	shares := float64(0)

	for i, record := range priceRecords {
		if i+1 < len(priceRecords) && record.Date.Month() != priceRecords[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		shares += amount / record.Price
		investmentMade = true

		records = append(records, models.WeeklyRecord{
			Time:     record.Date,
			NetWorth: CalculateNetWorth(record.Price, shares, 0),
			SnpValue: CalculateSnPValue(record.Price, shares, 0),
			Shares:   shares,
			Reserves: 0,
		})
	}

	return records
}

// Strat 2: Value Averaging
func VA(targetGrowth float64, r models.Records) []models.WeeklyRecord {
	reserves := 0.0
	shares := 0.0
	records := []models.WeeklyRecord{}
	investmentMade := false
	monthsCount := 1

	for i, record := range r {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		reserves += 1000
		snpValue := shares * record.Price
		targetSnPValue := float64(monthsCount) * targetGrowth
		amountToMove := targetSnPValue - snpValue
		monthsCount++

		if amountToMove <= reserves {
			reserves -= amountToMove
			shares += amountToMove / record.Price
		} else {
			shares += reserves / record.Price
			reserves = 0
		}

		records = append(records, models.WeeklyRecord{
			Time:     record.Date,
			NetWorth: CalculateNetWorth(record.Price, shares, reserves),
			SnpValue: CalculateSnPValue(record.Price, shares, reserves),
			Shares:   shares,
			Reserves: reserves,
		})

		investmentMade = true
	}

	return records
}

type DynamicVAConfig struct {
	BottomRatio          float64
	TopRatio             float64
	ReducingMultiplier   float64
	IncreasingMultiplier float64
}

// Strat 3: Dynamic Value Averaging
func DynamicVA(targetGrowth float64, r models.Records, cfg DynamicVAConfig) []models.WeeklyRecord {
	reserves := 0.0
	shares := 0.0
	records := []models.WeeklyRecord{}
	targetSnPValue := 0.0
	investmentMade := false

	for i, record := range r {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		reserves += 1000
		snpValue := shares * record.Price
		targetSnPValue += targetGrowth
		amountToMove := targetSnPValue - snpValue

		if amountToMove <= reserves {
			reserves -= amountToMove
			shares += amountToMove / record.Price
		} else {
			shares += reserves / record.Price
			reserves = 0
		}

		cashBufferRatio := reserves / targetGrowth
		if cashBufferRatio < cfg.BottomRatio {
			targetGrowth = max(targetGrowth*cfg.ReducingMultiplier, 1000)
		} else if cashBufferRatio > cfg.TopRatio {
			targetGrowth *= cfg.IncreasingMultiplier
		}

		records = append(records, models.WeeklyRecord{
			Time:     record.Date,
			NetWorth: CalculateNetWorth(record.Price, shares, reserves),
			SnpValue: CalculateSnPValue(record.Price, shares, reserves),
			Shares:   shares,
			Reserves: reserves,
		})

		investmentMade = true
	}

	return records
}

// calculateAverage computes the average of a slice of float64 numbers
func calculateAverage(prices []float64) float64 {
	sum := 0.0
	for _, price := range prices {
		sum += price
	}
	return sum / float64(len(prices))
}

// Strat 4: Buy Low, Sell High
func BuyLowSellHigh(r models.Records) []models.WeeklyRecord {
	reserves := 0.0
	shares := 0.0
	records := []models.WeeklyRecord{}
	windowSize := 20
	buyThreshold := 1.00
	sellThreshold := 1.03
	investmentMade := false

	var recentPrices []float64

	for i, record := range r {
		recentPrices = append(recentPrices, record.Price)
		if len(recentPrices) > windowSize {
			recentPrices = recentPrices[1:]
		}

		if len(recentPrices) < windowSize {
			continue
		}

		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		investmentMade = true

		avgPrice := calculateAverage(recentPrices)
		reserves += 1000

		if record.Price < buyThreshold*avgPrice {
			sharesToBuy := reserves / record.Price
			shares += sharesToBuy
			reserves -= sharesToBuy * record.Price
		}

		if record.Price > sellThreshold*avgPrice {
			proceeds := shares * record.Price
			reserves += proceeds
			shares = 0
		}

		records = append(records, models.WeeklyRecord{
			Time:     record.Date,
			NetWorth: CalculateNetWorth(record.Price, shares, reserves),
			SnpValue: CalculateSnPValue(record.Price, shares, reserves),
			Shares:   shares,
			Reserves: reserves,
		})
	}

	return records
}

// Strat 5: Mattress
func Mattress(r models.Records) []models.WeeklyRecord {
	reserves := 0.0
	records := []models.WeeklyRecord{}
	investmentMade := false

	for i, record := range r {
		if i+1 < len(r) && record.Date.Month() != r[i+1].Date.Month() {
			investmentMade = false
		}
		if investmentMade {
			continue
		}

		reserves += 1000
		investmentMade = true

		records = append(records, models.WeeklyRecord{
			Time:     record.Date,
			NetWorth: reserves,
			SnpValue: 0,
			Shares:   0,
			Reserves: reserves,
		})
	}

	return records
}

// Prints it out using tabwriter
func CompareStrats(out io.Writer, priceRecords models.Records, results []models.WeeklyRecords) {
	fmt.Fprintln(out, "Number of Strategies:", len(results))
	fmt.Fprintln(out)

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Strategy\tNet Worth\tSnP Value\tNo of Shares\tReserves Value\t")

	for _, records := range results {
		fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%.2f\t%.2f\t\n",
			records.Strategy,
			records.Records[len(records.Records)-1].NetWorth,
			records.Records[len(records.Records)-1].SnpValue,
			records.Records[len(records.Records)-1].Shares,
			records.Records[len(records.Records)-1].Reserves,
		)
	}
	fmt.Fprintln(w)

	w.Flush()
}
