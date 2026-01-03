package finance

import (
	"fmt"
	"time"
)

type RateChange struct {
	Date time.Time
	Rate float64
}

func RunTrackerMode(details LoanDetails, interestPaidYTD float64, rateHistory []RateChange) {
	fmt.Println("=== MODE 2: Loan Tracker ===")

	now := time.Now()

	endDate := details.StartDate.AddDate(0, details.LoanLengthMonths, 0)

	getRateAtDate := func(dateInQuestion time.Time) float64 {
		currentRate := details.InterestRate
		for _, change := range rateHistory {
			if change.Date.Before(dateInQuestion) || change.Date.Equal(dateInQuestion) {
				currentRate = change.Rate
			}
		}
		return currentRate
	}

	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	trackFrom := startOfYear
	if details.StartDate.After(startOfYear) {
		trackFrom = details.StartDate
	}

	var expectedInterestYTD float64

	monthsToCalculate := 0
	if now.After(trackFrom) {
		monthsToCalculate = int(now.Month()) - int(trackFrom.Month())
		if now.Year() > trackFrom.Year() {
			monthsToCalculate = int(now.Month())
		}
	}

	monthlyCorePayment := details.TotalLoanAmount / float64(details.LoanLengthMonths)

	for i := 0; i < monthsToCalculate; i++ {
		calcDate := trackFrom.AddDate(0, i, 0)

		rateForMonth := getRateAtDate(calcDate)

		totalMonthsElapsed := CalculateMonthsElapsedAtDate(details.StartDate, calcDate)
		coreRepaid := monthlyCorePayment * float64(totalMonthsElapsed)
		outstanding := details.TotalLoanAmount - coreRepaid
		if outstanding < 0 {
			outstanding = 0
		}

		monthlyInt := outstanding * (rateForMonth / 100) / 12.0
		expectedInterestYTD += monthlyInt
	}

	currentRate := getRateAtDate(now)
	totalMonthsElapsed := CalculateMonthsElapsedAtDate(details.StartDate, now)
	coreRepaidTotal := monthlyCorePayment * float64(totalMonthsElapsed)
	outstandingCoreNow := details.TotalLoanAmount - coreRepaidTotal
	if outstandingCoreNow < 0 {
		outstandingCoreNow = 0
	}

	currentMonthlyCost := outstandingCoreNow * (currentRate / 100) / 12.0

	monthsLeftInYear := 12 - int(now.Month()) + 1

	var liabilityForRestOfYear float64
	tempCore := outstandingCoreNow

	for i := 0; i < monthsLeftInYear; i++ {
		futureDate := now.AddDate(0, i, 0)
		if futureDate.After(endDate) {
			break
		}

		mInt := tempCore * (currentRate / 100) / 12.0
		liabilityForRestOfYear += mInt

		tempCore -= monthlyCorePayment
	}

	balance := expectedInterestYTD - interestPaidYTD
	tolerance := 0.05

	PrintSeparator()
	fmt.Printf("Loan start date:		%s\n", details.StartDate.Format("2006-01-02"))
	fmt.Printf("Loan end date:			%s\n", endDate.Format("2006-01-02"))
	PrintSeparator()
	fmt.Printf("Loan Progress:			%d / %d months\n", totalMonthsElapsed, details.LoanLengthMonths)
	fmt.Printf("Outstanding Core:		%s\n", FormatCurrency(outstandingCoreNow))
	fmt.Printf("Current Interest Rate:		%.2f%%\n", currentRate)
	PrintSeparator()
	fmt.Printf("Tracking Period:		%s to Now (%d billed months)\n", trackFrom.Format("Jan 2006"), monthsToCalculate)
	fmt.Printf("Expected Interest YTD:		%s\n", FormatCurrency(expectedInterestYTD))
	fmt.Printf("Interest Paid YTD:		%s\n", FormatCurrency(interestPaidYTD))
	PrintSeparator()

	if balance > tolerance {
		fmt.Println("STATUS: [RUNNING BEHIND]")
		fmt.Printf("Arrears (YTD):				%s\n", FormatCurrency(balance))

		totalLumpSum := balance + liabilityForRestOfYear

		fmt.Println("\n--- Correction Plan (Remainder of Year) ---")
		fmt.Printf("To clear arrears and cover the rest of %d:\n", now.Year())
		fmt.Printf("- pay a lump sum of %s by December 31st\n", FormatCurrency(totalLumpSum))

		if monthsLeftInYear > 0 {
			monthlyCatchup := totalLumpSum / float64(monthsLeftInYear)
			fmt.Printf("- alternatively increase monthly interest payments to %s\n", FormatCurrency(monthlyCatchup))
		}
	} else {
		fmt.Println("STATUS: [ON TRACK]")
		fmt.Printf("You are covering the YTD interest requirements!\n")
		fmt.Printf("Projected new monthly cost:	%s\n", FormatCurrency(currentMonthlyCost))

		surplus := -balance
		if surplus > tolerance {
			fmt.Printf("You have a YTD surplus of:	%s\n", FormatCurrency(surplus))
		}
	}
	PrintSeparator()
}

//func RunTrackerMode(details LoanDetails, interestPaidSoFar float64) {
//	fmt.Println("=== MODE 2: Loan Tracker ===")
//
//	monthsElapsed := CalculateMonthsElapsed(details.StartDate)
//	if monthsElapsed > details.LoanLengthMonths {
//		monthsElapsed = details.LoanLengthMonths
//	}
//	monthsRemaining := details.LoanLengthMonths - monthsElapsed
//	endDate := details.StartDate.AddDate(0, details.LoanLengthMonths, 0)
//
//	monthlyCorePayment := details.TotalLoanAmount / float64(details.LoanLengthMonths)
//	coreRepaidSoFar := monthlyCorePayment * float64(monthsElapsed)
//	outstandingCore := details.TotalLoanAmount - coreRepaidSoFar
//
//	if outstandingCore < 0 {
//		outstandingCore = 0
//	}
//
//	currentAnnualInterest := outstandingCore * (details.InterestRate / 100.0)
//	currentMonthlyInterestCost := currentAnnualInterest / 12.0
//
//	var projectedRemainingInterest float64
//	tempCore := outstandingCore
//	for i := 0; i < monthsRemaining; i++ {
//		monthlyInt := tempCore * (details.InterestRate / 100.0) / 12.0
//		projectedRemainingInterest += monthlyInt
//		tempCore -= monthlyCorePayment
//	}
//
//	now := time.Now()
//
//	var totalExpectedToDate float64
//	p := details.TotalLoanAmount
//	for i := 0; i < monthsElapsed; i++ {
//		mInt := p * (details.InterestRate / 100.0) / 12.0
//		totalExpectedToDate += mInt
//		p -= monthlyCorePayment
//	}
//
//	balance := totalExpectedToDate - interestPaidSoFar
//
//	PrintSeparator()
//	fmt.Printf("Loan start date:            %s\n", details.StartDate.Format("2006-01-02"))
//	fmt.Printf("Loan end date:                 %s\n", endDate.Format("2006-01-02"))
//	PrintSeparator()
//	fmt.Printf("Loan progress:                 %d / %d months\n", monthsElapsed, details.LoanLengthMonths)
//	fmt.Printf("Core repaid so far:             £%.2f\n", coreRepaidSoFar)
//	fmt.Printf("Outstanding core:           £%.2f\n", outstandingCore)
//	PrintSeparator()
//	fmt.Printf("Current interest rate:           %.2f%%\n", details.InterestRate)
//	fmt.Printf("New monthly interest:        £%.2f (based on outstanding balance)\n", currentMonthlyInterestCost)
//	fmt.Printf("Projected remaining interest:  £%.2f\n", projectedRemainingInterest)
//	PrintSeparator()
//
//	monthsLeftInYear := 12 - int(now.Month()) + 1
//
//	var liabilityForRestOfYear float64
//	tempCoreForYear := outstandingCore
//	for i := 0; i < monthsLeftInYear; i++ {
//		if i >= monthsRemaining {
//			break
//		}
//
//		mInt := tempCoreForYear * (details.InterestRate / 100.0) / 12.0
//		liabilityForRestOfYear += mInt
//		tempCoreForYear -= monthlyCorePayment
//	}
//
//	tolerance := 0.05
//
//	if balance > tolerance {
//		fmt.Println("STATUS: [RUNNING BEHIND]")
//		fmt.Printf("Estimated arrears based on current rate: £%.2f\n", balance)
//
//		totalLumpSum := balance + liabilityForRestOfYear
//
//		fmt.Println("\n--- Correction plan ---")
//		fmt.Printf("To clear arrears and cover the rest of %d:\n", now.Year())
//		fmt.Printf("- Pay a total lump sum of £%.2f by December 31st", totalLumpSum)
//
//		if monthsLeftInYear > 0 {
//			monthlyCatchup := totalLumpSum / float64(monthsLeftInYear)
//			fmt.Printf("- Alternatively increase monthly payments to £%.2f\n", monthlyCatchup)
//		} else {
//			fmt.Println("(No months left in year for catchup payments")
//		}
//	} else {
//		fmt.Println("STATUS [ON TRACK]")
//		fmt.Println("Looks like you're covering the interest requirements!")
//		fmt.Printf("Expected monthly interest cost: £%.2f\n", currentMonthlyInterestCost)
//
//		surplus := -balance
//		if surplus > tolerance {
//			fmt.Printf("(You have a surplus of £%.2f)\n", surplus)
//		}
//	}
//}
