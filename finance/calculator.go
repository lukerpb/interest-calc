package finance

import (
	"math"
	"time"
)

type LoanDetails struct {
	TotalLoanAmount  float64
	InterestRate     float64
	LoanLengthMonths int
	StartDate        time.Time
}

func RoundTo2DP(val float64) float64 {
	return math.Round(val*100) / 100
}

func CalculateTotalInterestFlat(details LoanDetails) float64 {
	years := float64(details.LoanLengthMonths) / 12.0
	interest := details.TotalLoanAmount * (details.InterestRate / 100.0) * years
	return RoundTo2DP(interest)
}

func CalculateMonthsElapsed(start time.Time) int {
	now := time.Now()
	if now.Before(start) {
		return 0
	}

	years := now.Year() - start.Year()
	months := int(now.Month()) - int(start.Month())

	totalMonths := (years * 12) + months
	return totalMonths
}
