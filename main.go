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
	TotalLoanAmount  float64            `json:"totalLoanAmount"`
	InterestRate     float64            `json:"interestRate"`
	LoanLengthMonths int                `json:"loanLengthMonths"`
	StartDate        string             `json:"startDate"`
	RateChanges      []RateChangeConfig `json:"rateChanges"`
}

type RateChangeConfig struct {
	WhenRateChanged string  `json:"whenRateChanged"`
	RateChangedTo   float64 `json:"rateChangedTo"`
}

func main() {
	mode := flag.String("mode", "plan", "Mode to run: 'plan' or 'track'")
	amount := flag.Float64("amount", 0, "Total loan amount")
	rate := flag.Float64("rate", 0, "Interest rate in %")
	months := flag.Int("months", 0, "Loan length in months")
	start := flag.String("start", "", "Start date (YYYY-MM-DD)")
	paid := flag.Float64("paid", 0, "Interest repaid year-to-date (track mode only)")
	newRate := flag.Float64("newRate", 0, "New interest rate in % effective from today")

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
		var history []finance.RateChange
		now := time.Now()
		loanEndDate := startDate.AddDate(0, finalMonths, 0)

		if now.Before(startDate) {
			fmt.Printf("Error: current date (%s) is before the loan start date (%s).\nTracking is not available yet.\n",
				now.Format("2006-01-02"), startDate.Format("2006-01-02"))
			os.Exit(1)
		}

		if now.After(loanEndDate) {
			fmt.Printf("Error: current date (%s) is after the loan end date (%s).\nTracking is no longer available.\n",
				now.Format("2006-01-02"), loanEndDate.Format("2006-01-02"))
			os.Exit(1)
		}

		for _, rateChange := range cfg.RateChanges {
			date, err := time.Parse("2006-01-02", rateChange.WhenRateChanged)
			if err != nil {
				fmt.Printf("Error parsing date '%s': %v\n", rateChange.WhenRateChanged, err)
				os.Exit(1)
			}

			if date.Before(startDate) {
				fmt.Printf("Error: rate change %s is before loan start date\n", rateChange.WhenRateChanged)
				os.Exit(1)
			}

			if date.After(loanEndDate) {
				fmt.Printf("Error: rate change %s is after loan end date\n", rateChange.WhenRateChanged)
				os.Exit(1)
			}

			if date.After(now) {
				fmt.Printf("Error: rate change %s is in the future - update config only when rate changes\n", rateChange.WhenRateChanged)
				os.Exit(1)
			}

			history = append(history, finance.RateChange{
				Date: date,
				Rate: rateChange.RateChangedTo,
			})
		}

		if *newRate > 0 {
			fmt.Printf("Note: applying new rate of %.2f%% effective today\n", *newRate)
			history = append(history, finance.RateChange{
				Date: now,
				Rate: *newRate,
			})
		}

		finance.RunTrackerMode(details, *paid, history)

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
