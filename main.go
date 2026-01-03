package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"interest-calc/finance"
	"os"
	"time"
)

type Config struct {
	TotalLoanAmount  float64 `json:"totalLoanAmount"`
	InterestRate     float64 `json:"interestRate"`
	LoanLengthMonths int     `json:"loanLengthMonths"`
	StartDate        string  `json:"startDate"`
}

func main() {
	mode := flag.String("mode", "plan", "Mode to run: 'plan' or 'track'")
	amount := flag.Float64("amount", 0, "Total loan amount")
	rate := flag.Float64("rate", 0, "Interest rate in %")
	months := flag.Int("months", 0, "Loan length in months")
	start := flag.String("start", "", "Start date (YYYY-MM-DD)")
	paid := flag.Float64("paid", 0, "Interest repaid so far (track mode only")

	flag.Parse()

	cfg := loadConfig("config.json")

	finalAmount := resolveFloat(*amount, cfg.TotalLoanAmount)
	finalRate := resolveFloat(*rate, cfg.InterestRate)
	finalMonths := resolveInt(*months, cfg.LoanLengthMonths)
	finalStart := resolveString(*start, cfg.StartDate)

	if finalAmount == 0 || finalMonths == 0 || finalStart == "" {
		fmt.Println("Error: missing loan details. Please provide via flags or config.json")
		flag.PrintDefaults()
		os.Exit(1)
	}

	startDate, err := time.Parse("2006-01-02", finalStart)
	if err != nil {
		fmt.Printf("Error parsing start date '%s': %v \nUse format YYYY-MM-DD\n", finalStart, err)
		os.Exit(1)
	}

	details := finance.LoanDetails{
		TotalLoanAmount:  finalAmount,
		InterestRate:     finalRate,
		LoanLengthMonths: finalMonths,
		StartDate:        startDate,
	}

	switch *mode {
	case "plan":
		finance.RunPlannerMode(details)
	case "track":
		finance.RunTrackerMode(details, *paid)
	default:
		fmt.Println("Error: unknown mode. Use 'plan' or 'track'.")
		os.Exit(1)
	}
}

func loadConfig(filename string) Config {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Warning: could not parse config file - checking flags...")
	}
	return config
}

func resolveFloat(flagVal, configVal float64) float64 {
	if flagVal != 0 {
		return flagVal
	}
	return configVal
}

func resolveInt(flagVal, configVal int) int {
	if flagVal != 0 {
		return flagVal
	}
	return configVal
}

func resolveString(flagVal, configVal string) string {
	if flagVal != "" {
		return flagVal
	}
	return configVal
}
