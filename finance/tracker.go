package finance

import (
	"fmt"
	"time"
)

type RateChange struct {
	Date time.Time
	Rate float64
}

const (
	StartColourGreen = "\033[32m"
	StartColourAmber = "\033[38;5;208m"
	StopColouring    = "\033[0m"
)

func RunTrackerMode(details LoanDetails, interestPaidYTD float64, rateHistory []RateChange) {
	fmt.Println("=== MODE 2: Loan Tracker ===")

	now := time.Now()

	getRateAtDate := func(dateInQuestion time.Time) float64 {
		currentRate := details.InterestRate
		for _, change := range rateHistory {
			if change.Date.Before(dateInQuestion) || change.Date.Equal(dateInQuestion) {
				currentRate = change.Rate
			}
		}
		return currentRate
	}

	totalMonthsElapsed := CalculateMonthsElapsedAtDate(details.StartDate, now)

	if totalMonthsElapsed > details.LoanLengthMonths {
		totalMonthsElapsed = details.LoanLengthMonths
	}

	percentageProgress := (float64(totalMonthsElapsed) / float64(details.LoanLengthMonths)) * 100.0

	monthlyCorePayment := details.TotalLoanAmount / float64(details.LoanLengthMonths)

	coreRepaidTotal := monthlyCorePayment * float64(totalMonthsElapsed)

	outstandingCoreNow := details.TotalLoanAmount - coreRepaidTotal
	if outstandingCoreNow < 0 {
		outstandingCoreNow = 0
	}

	startOfCurrentYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	var interestPaidPreviousYears float64

	if details.StartDate.Before(startOfCurrentYear) {
		monthsPrior := CalculateMonthsElapsedAtDate(details.StartDate, startOfCurrentYear)

		for i := 0; i < monthsPrior; i++ {
			calcDate := details.StartDate.AddDate(0, i, 0)

			rate := getRateAtDate(calcDate)

			cRepaid := monthlyCorePayment * float64(i)
			cOutstanding := details.TotalLoanAmount - cRepaid

			mInt := cOutstanding * (rate / 100.0) / 12.0
			interestPaidPreviousYears += mInt
		}
	}

	totalInterestPaidLifetime := interestPaidPreviousYears + interestPaidYTD
	totalPaidLifetime := coreRepaidTotal + totalInterestPaidLifetime

	var expectedInterestYTD float64

	trackFrom := startOfCurrentYear
	if details.StartDate.After(startOfCurrentYear) {
		trackFrom = details.StartDate
	}

	monthsToCalculateYTD := 0
	if now.After(trackFrom) {
		monthsToCalculateYTD = int(now.Month()) - int(trackFrom.Month())
		if now.Year() > trackFrom.Year() {
			monthsToCalculateYTD = int(now.Month()) // Handle Dec -> Jan crossover edge case
		}
	}

	for i := 0; i < monthsToCalculateYTD; i++ {
		calcDate := trackFrom.AddDate(0, i, 0)
		rateForMonth := getRateAtDate(calcDate)

		monthsFromStart := CalculateMonthsElapsedAtDate(details.StartDate, calcDate)

		coreRepaid := monthlyCorePayment * float64(monthsFromStart)
		outstanding := details.TotalLoanAmount - coreRepaid

		monthlyInt := outstanding * (rateForMonth / 100.0) / 12.0
		expectedInterestYTD += monthlyInt
	}

	currentRate := getRateAtDate(now)
	monthsRemaining := details.LoanLengthMonths - totalMonthsElapsed

	var projectedRemainingInterest float64
	tempCore := outstandingCoreNow

	for i := 0; i < monthsRemaining; i++ {
		mInt := tempCore * (currentRate / 100.0) / 12.0
		projectedRemainingInterest += mInt
		tempCore -= monthlyCorePayment
	}

	totalOutstanding := outstandingCoreNow + projectedRemainingInterest

	balance := expectedInterestYTD - interestPaidYTD

	monthsLeftInYear := 12 - int(now.Month()) + 1
	var liabilityForRestOfYear float64
	tempCoreYear := outstandingCoreNow
	for i := 0; i < monthsLeftInYear; i++ {
		if i >= monthsRemaining {
			break
		}
		mInt := tempCoreYear * (currentRate / 100.0) / 12.0
		liabilityForRestOfYear += mInt
		tempCoreYear -= monthlyCorePayment
	}

	PrintSeparator()
	fmt.Printf(" Loan Progress:			%d / %d months (%.1f%%)\n", totalMonthsElapsed, details.LoanLengthMonths, percentageProgress)
	PrintSeparator()
	fmt.Println(" PAID SO FAR (Lifetime):")
	fmt.Printf(" - Total Paid:			%s\n", FormatCurrency(totalPaidLifetime))
	fmt.Printf("   - Core:			%s\n", FormatCurrency(coreRepaidTotal))
	fmt.Printf("   - Interest:			%s\n", FormatCurrency(totalInterestPaidLifetime))
	PrintSeparator()
	fmt.Println(" REMAINING PREDICTION (Assumes current rate holds):")
	fmt.Printf(" - Total Outstanding:		%s\n", FormatCurrency(totalOutstanding))
	fmt.Printf("   - Core:			%s\n", FormatCurrency(outstandingCoreNow))
	fmt.Printf("   - Interest:			%s\n", FormatCurrency(projectedRemainingInterest))
	PrintSeparator()
	fmt.Printf(" Current Interest Rate:		%.2f%%\n", currentRate)
	PrintSeparator()
	fmt.Printf(" Tracking Period:		%s to Now (%d billed months)\n", trackFrom.Format("Jan 2006"), monthsToCalculateYTD)
	fmt.Printf(" Expected Interest YTD:		%s\n", FormatCurrency(expectedInterestYTD))
	fmt.Printf(" Interest Paid YTD:		%s\n", FormatCurrency(interestPaidYTD))
	PrintSeparator()

	tolerance := 0.05
	if balance > tolerance {
		fmt.Println(" STATUS: " + MakeThisAmber("[RUNNING BEHIND]"))
		fmt.Printf(" Arrears (YTD):%s\n", FormatCurrency(balance))

		totalLumpSum := balance + liabilityForRestOfYear

		fmt.Println("\n--- Correction Plan (Remainder of Year) ---")
		fmt.Printf("To clear arrears and cover the rest of %d:\n", now.Year())
		fmt.Printf("- pay a lump sum of %s by December 31st\n", FormatCurrency(totalLumpSum))

		if monthsLeftInYear > 0 {
			monthlyCatchup := totalLumpSum / float64(monthsLeftInYear)
			fmt.Printf("- alternatively increase monthly interest payments to %s\n", FormatCurrency(monthlyCatchup))
		}
	} else {
		fmt.Println(" STATUS: " + MakeThisGreen("[ON TRACK]"))
		fmt.Println(" You are covering the YTD interest requirements!")

		currentMonthlyCost := outstandingCoreNow * (currentRate / 100.0) / 12.0
		fmt.Printf(" Projected Monthly Cost:	%s\n", FormatCurrency(currentMonthlyCost))

		surplus := -balance
		if surplus > tolerance {
			fmt.Printf(" (You have a YTD surplus of:%s)\n", FormatCurrency(surplus))
		}
	}
	PrintSeparator()
}

func MakeThisGreen(text string) string {
	return StartColourGreen + text + StopColouring
}

func MakeThisAmber(text string) string {
	return StartColourAmber + text + StopColouring
}
