package finance

import "fmt"

func RunPlannerMode(details LoanDetails) {
	fmt.Println("=== MODE 1: Loan Planner ===")

	totalInterest := CalculateTotalInterestFlat(details)
	totalPayable := details.TotalLoanAmount + totalInterest
	endDate := details.StartDate.AddDate(0, details.LoanLengthMonths, 0)

	monthlyCore := RoundTo2DP(details.TotalLoanAmount / float64(details.LoanLengthMonths))
	monthlyInterest := RoundTo2DP(totalInterest / float64(details.LoanLengthMonths))
	standardMonthlyTotal := monthlyCore + monthlyInterest

	totalScheduledPayment := standardMonthlyTotal + float64(details.LoanLengthMonths)
	diff := RoundTo2DP(totalPayable - totalScheduledPayment)
	lastMonthPayment := standardMonthlyTotal + diff

	PrintSeparator()
	fmt.Printf("Loan start date: %s\n", details.StartDate.Format("2006-01-02"))
	fmt.Printf("Loan end date: %s\n", endDate.Format("2006-01-02"))
	PrintSeparator()
	fmt.Printf("Total core:		£%.2f\n", details.TotalLoanAmount)
	fmt.Printf("Total interest:	£%.2f\n", totalInterest)
	fmt.Printf("Total payable: 	£%.2f\n", totalPayable)
	PrintSeparator()
	fmt.Printf("Monthly core contribution: 		£%.2f\n", monthlyCore)
	fmt.Printf("Monthly interest contribution: 	£%.2f\n", monthlyInterest)
	fmt.Printf("Standard monthly payment: 		£%.2f\n", standardMonthlyTotal)

	if diff != 0 {
		fmt.Printf("\n* Adjustment required for final month:\n")
		fmt.Printf("	Final payment: 			£%.2f\n", lastMonthPayment)
	} else {
		fmt.Printf("\n* No adjustment required for final month:\n")
	}
}

func PrintSeparator() {
	fmt.Println("------------------------------------------------")
}
