package finance

import (
	"fmt"
	"math"
)

func RunPlannerMode(details LoanDetails) {
	fmt.Println("=== MODE 1: Loan Planner ===")

	totalInterest := CalculateTotalCompoundInterest(details)
	totalPayable := details.TotalLoanAmount + totalInterest
	endDate := details.StartDate.AddDate(0, details.LoanLengthMonths, 0)

	exactMonthlyPayment := totalPayable / float64(details.LoanLengthMonths)
	standardMonthlyTotal := math.Ceil(exactMonthlyPayment*100) / 100

	paymentsExcludingLastOne := standardMonthlyTotal * float64(details.LoanLengthMonths-1)
	lastMonthPayment := totalPayable - paymentsExcludingLastOne
	lastMonthPayment = math.Round(lastMonthPayment*100) / 100

	monthlyCore := math.Round((details.TotalLoanAmount/float64(details.LoanLengthMonths))*100) / 100
	monthlyInterest := standardMonthlyTotal - monthlyCore

	diff := standardMonthlyTotal - lastMonthPayment
	hasDiff := diff > 0.001

	PrintSeparator()
	fmt.Printf("Loan start date:		%s\n", details.StartDate.Format("2006-01-02"))
	fmt.Printf("Loan end date:			%s\n", endDate.Format("2006-01-02"))
	PrintSeparator()
	fmt.Printf("Total core:			%s\n", FormatCurrency(details.TotalLoanAmount))
	fmt.Printf("Total interest:			%s\n", FormatCurrency(totalInterest))
	fmt.Printf("Total payable:			%s\n", FormatCurrency(totalPayable))
	PrintSeparator()
	fmt.Printf("Monthly core contribution:	%s\n", FormatCurrency(monthlyCore))
	fmt.Printf("Monthly interest contribution:	%s\n", FormatCurrency(monthlyInterest))
	PrintSeparator()
	fmt.Println("PAYMENT SCHEDULE:")

	if !hasDiff {
		fmt.Printf("%d monthly payments of:	%s\n", details.LoanLengthMonths, FormatCurrency(standardMonthlyTotal))
	} else {
		fmt.Printf(" - %d monthly payments of:	%s\n", details.LoanLengthMonths-1, FormatCurrency(standardMonthlyTotal))
		fmt.Printf(" - Final payment of:		%s\n", FormatCurrency(lastMonthPayment))

		fmt.Println("	|") // Visual connector

		finalMonthCore := monthlyCore
		finalMonthInterest := lastMonthPayment - finalMonthCore

		if lastMonthPayment < monthlyCore {
			finalMonthCore = lastMonthPayment
			finalMonthInterest = 0
		}

		fmt.Printf("	L Final Month Breakdown:\n")
		fmt.Printf("  	- Core:			%s\n", FormatCurrency(finalMonthCore))
		if finalMonthInterest > 0 {
			fmt.Printf("  	- Interest:		%s\n", FormatCurrency(finalMonthInterest))
		} else {
			fmt.Printf("  	- Interest:              £       0.00\n")
		}
	}
}

func PrintSeparator() {
	fmt.Println("------------------------------------------------")
}
