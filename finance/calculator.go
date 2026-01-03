package finance

import (
	"fmt"
	"math"
	"strings"
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

func CalculateTotalCompoundInterest(details LoanDetails) float64 {
	years := float64(details.LoanLengthMonths) / 12.0
	rateAsDecimal := details.InterestRate / 100.0

	// Calculate total amount using the annual compound interest formula: A = P(1 + r)^t
	totalRepayable := details.TotalLoanAmount * math.Pow(1+rateAsDecimal, years)
	interest := totalRepayable - details.TotalLoanAmount

	return RoundTo2DP(interest)
}

func CalculateMonthsElapsedAtDate(start, target time.Time) int {
	if target.Before(start) {
		return 0
	}
	years := target.Year() - start.Year()
	months := int(target.Month()) - int(start.Month())
	return (years * 12) + months
}

func FormatCurrency(amount float64) string {
	rounded := RoundTo2DP(amount)
	pounds := int64(rounded)
	pennies := int64(math.Round((rounded - float64(pounds)) * 100))
	if pennies < 0 {
		pennies = -pennies
	}

	poundString := fmt.Sprintf("%d", pounds)
	var result []string
	// Reverse the string to easily insert thousand separators every 3 digits
	for i, c := range reverse(poundString) {
		if i > 0 && i%3 == 0 {
			result = append(result, ",")
		}
		result = append(result, c)
	}
	// Re-reverse to restore the original order with commas included
	withCommas := strings.Join(reverse(strings.Join(result, "")), "")

	// Pad to ensure numerical alignment in console output
	paddedPoundAmount := fmt.Sprintf("%9s", withCommas)
	return fmt.Sprintf("£%s.%02d", paddedPoundAmount, pennies)
}

func reverse(s string) []string {
	var r []string
	for _, c := range s {
		r = append([]string{string(c)}, r...)
	}
	return r
}
