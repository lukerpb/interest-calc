# Interest Calculator

A Go-based command-line tool designed to help you plan and track loan interest. It supports annual compound interest calculations and provides detailed tracking of interest payments over time, including support for historical interest rate fluctuations.

## Features

- **Loan Planning**: Calculate total interest, monthly payments, and see a full payment schedule.
- **Interest Tracking**: Monitor your progress, calculate lifetime interest paid, and verify if your year-to-date (YTD) payments match expected interest accrual.
- **Historical Rate Support**: Account for multiple interest rate changes throughout the life of the loan.
- **Correction Plans**: If you are behind on interest payments, the tool suggests a lump-sum payment or a monthly catch-up amount to get back on track.

## Modes of Operation

The calculator runs in two primary modes, controlled by the `-mode` flag.

### 1. Planner Mode (`-mode=plan`)
Use this mode when you want to see the long-term cost of a loan. It calculates:
- Total interest over the life of the loan using the compound interest formula $A = P(1 + r)^t$.
- The exact monthly payment required to clear the loan (core + interest).
- A breakdown of the final payment, which may differ slightly due to rounding.

### 2. Tracker Mode (`-mode=track`)
Use this mode to monitor an active loan. It provides:
- **Progress Overview**: How many months have elapsed and what percentage of the loan term is complete.
- **Lifetime Totals**: Total amount paid so far, broken down by core repayment and interest.
- **YTD Analysis**: Compares the interest you've actually paid this calendar year (provided via `-paid`) against the expected interest based on the loan's history.
- **Projections**: Predicts remaining interest and total outstanding balance based on the current rate.

---

## Configuration

You can manage your loan details via a `config.json` file kept at the base of the project.

### Example file
```json
{
  "totalLoanAmount": 25000,
  "interestRate": 4.5,
  "loanLengthMonths": 48,
  "startDate": "2023-05-15",
  "rateChanges": [
    {
      "whenRateChanged": "2023-11-01",
      "rateChangedTo": 5.25
    },
    {
      "whenRateChanged": "2024-02-01",
      "rateChangedTo": 5.0
    }
  ]
}
```

- `totalLoanAmount`: The initial principal of the loan.
- `interestRate`: The base annual interest rate (in %).
- `loanLengthMonths`: The duration of the loan in months.
- `startDate`: When the loan began (format: `YYYY-MM-DD`).
- `rateChanges`: An optional list of historical rate changes. Each entry requires a date (`whenRateChanged`) and the new rate (`rateChangedTo`).

---

## Command-Line Flags

Flags take precedence over values defined in `config.json`.

| Flag | Mode | Description |
|------|------|-------------|
| `-mode` | Global | Choose between `plan` or `track`. Default is `plan`. |
| `-amount` | Global | Total loan amount (e.g., `15000`). |
| `-rate` | Global | Annual interest rate in % (e.g., `3.5`). |
| `-months` | Global | Total loan duration in months (e.g., `36`). |
| `-start` | Global | Loan start date in `YYYY-MM-DD` format. |
| `-paid` | Track | The total interest you have repaid **this calendar year**. |
| `-newRate` | Track | Apply a new interest rate effective from **today** (useful for "what-if" scenarios). |

---

## How to Run

### Installation
Ensure you have [Go](https://go.dev/) installed on your system.

1. Clone or download the repository.
2. Navigate to the project directory:
   ```bash
   cd interest-calc
   ```

### Examples

#### 1. Running a simple plan via flags
If you don't have a config file, you must provide the core details via flags:
```bash
go run main.go -mode=plan -amount=10000 -rate=5 -months=24 -start=2024-01-01
```

#### 2. Running tracker using a config file
Assuming you have a `config.json` configured with your loan details:
```bash
go run main.go -mode=track -paid=120.50
```
This will check if your YTD interest payment of £120.50 is sufficient based on the rates defined in your config.

#### 3. Overriding config values
You can use the config for most values but override a specific one (like the interest rate) using flags:
```bash
go run main.go -mode=track -paid=300.00 -rate=6.5
```

#### 4. Projecting a new rate change
To see how a rate change today would affect your remaining balance:
```bash
go run main.go -mode=track -paid=250.00 -newRate=7.25
```
