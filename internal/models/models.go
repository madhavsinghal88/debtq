package models

import "time"

// ExpenseCategory represents expense categories
type ExpenseCategory string

const (
	CategoryFood          ExpenseCategory = "food"
	CategoryTransport     ExpenseCategory = "transport"
	CategoryEntertainment ExpenseCategory = "entertainment"
	CategoryUtilities     ExpenseCategory = "utilities"
	CategoryShopping      ExpenseCategory = "shopping"
	CategoryHealth        ExpenseCategory = "health"
	CategoryEducation     ExpenseCategory = "education"
	CategoryOther         ExpenseCategory = "other"
)

// Expense represents a single expense entry
type Expense struct {
	ID          string          `json:"id"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	Category    ExpenseCategory `json:"category"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
}

// TransactionType for borrowing/lending
type TransactionType string

const (
	Borrowed TransactionType = "borrowed"
	Lent     TransactionType = "lent"
)

// DebtTransaction represents money borrowed or lent
type DebtTransaction struct {
	ID          string          `json:"id"`
	Type        TransactionType `json:"type"`
	PersonName  string          `json:"person_name"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	IsSettled   bool            `json:"is_settled"`
	SettledDate *time.Time      `json:"settled_date,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

// InvestmentType represents types of investments
type InvestmentType string

const (
	InvestmentStocks      InvestmentType = "stocks"
	InvestmentMutualFunds InvestmentType = "mutual_funds"
	InvestmentGold        InvestmentType = "gold"
	InvestmentSilver      InvestmentType = "silver"
	InvestmentFD          InvestmentType = "fixed_deposit"
	InvestmentPPF         InvestmentType = "ppf"
	InvestmentCrypto      InvestmentType = "crypto"
	InvestmentRealEstate  InvestmentType = "real_estate"
	InvestmentOther       InvestmentType = "other"
)

// Investment represents an investment entry
type Investment struct {
	ID             string         `json:"id"`
	Type           InvestmentType `json:"type"`
	Name           string         `json:"name"`
	InvestedAmount float64        `json:"invested_amount"`
	CurrentValue   float64        `json:"current_value"`
	Units          float64        `json:"units,omitempty"`
	PurchaseDate   time.Time      `json:"purchase_date"`
	Notes          string         `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// SavingsTarget represents a savings goal
type SavingsTarget struct {
	ID            string    `json:"id"`
	ProductName   string    `json:"product_name"`
	TargetAmount  float64   `json:"target_amount"`
	CurrentAmount float64   `json:"current_amount"`
	TargetDate    time.Time `json:"target_date"`
	Description   string    `json:"description,omitempty"`
	IsCompleted   bool      `json:"is_completed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SavingsContribution represents a contribution towards a savings target
type SavingsContribution struct {
	ID        string    `json:"id"`
	TargetID  string    `json:"target_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Data holds all the application data
type Data struct {
	Expenses             []Expense             `json:"expenses"`
	DebtTransactions     []DebtTransaction     `json:"debt_transactions"`
	Investments          []Investment          `json:"investments"`
	SavingsTargets       []SavingsTarget       `json:"savings_targets"`
	SavingsContributions []SavingsContribution `json:"savings_contributions"`
}

// NetWorth calculates total net worth from investments
func (d *Data) NetWorth() float64 {
	var total float64
	for _, inv := range d.Investments {
		total += inv.CurrentValue
	}
	return total
}

// TotalBorrowed returns total amount borrowed (unsettled)
func (d *Data) TotalBorrowed() float64 {
	var total float64
	for _, dt := range d.DebtTransactions {
		if dt.Type == Borrowed && !dt.IsSettled {
			total += dt.Amount
		}
	}
	return total
}

// TotalLent returns total amount lent (unsettled)
func (d *Data) TotalLent() float64 {
	var total float64
	for _, dt := range d.DebtTransactions {
		if dt.Type == Lent && !dt.IsSettled {
			total += dt.Amount
		}
	}
	return total
}

// MonthlyExpenses returns total expenses for a given month
func (d *Data) MonthlyExpenses(year int, month time.Month) float64 {
	var total float64
	for _, exp := range d.Expenses {
		if exp.Date.Year() == year && exp.Date.Month() == month {
			total += exp.Amount
		}
	}
	return total
}

// GetSavingsProgress returns progress percentage for a savings target
func (st *SavingsTarget) GetProgress() float64 {
	if st.TargetAmount == 0 {
		return 100
	}
	return (st.CurrentAmount / st.TargetAmount) * 100
}

// DaysRemaining returns days until target date
func (st *SavingsTarget) DaysRemaining() int {
	return int(time.Until(st.TargetDate).Hours() / 24)
}

// RequiredMonthlySavings calculates how much needs to be saved per month
func (st *SavingsTarget) RequiredMonthlySavings() float64 {
	remaining := st.TargetAmount - st.CurrentAmount
	if remaining <= 0 {
		return 0
	}
	months := float64(st.DaysRemaining()) / 30.0
	if months <= 0 {
		return remaining
	}
	return remaining / months
}
