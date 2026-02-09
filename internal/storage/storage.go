package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/debtq/debtq/internal/config"
	"github.com/debtq/debtq/internal/models"
	"github.com/google/uuid"
)

// Storage handles data persistence
type Storage struct {
	config *config.Config
	data   *models.Data
}

// New creates a new storage instance
func New(cfg *config.Config) (*Storage, error) {
	s := &Storage{
		config: cfg,
		data:   &models.Data{},
	}

	if err := s.Load(); err != nil {
		// If file doesn't exist, initialize empty data
		if os.IsNotExist(err) {
			s.data = &models.Data{
				Expenses:             []models.Expense{},
				DebtTransactions:     []models.DebtTransaction{},
				Investments:          []models.Investment{},
				SavingsTargets:       []models.SavingsTarget{},
				SavingsContributions: []models.SavingsContribution{},
			}
			return s, nil
		}
		return nil, err
	}

	return s, nil
}

// Load loads data from file
func (s *Storage) Load() error {
	dataPath := s.config.DataFile

	// Ensure directory exists
	dir := filepath.Dir(dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.data)
}

// Save saves data to file
func (s *Storage) Save() error {
	dataPath := s.config.DataFile

	// Ensure directory exists
	dir := filepath.Dir(dataPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataPath, data, 0644)
}

// GetData returns the current data
func (s *Storage) GetData() *models.Data {
	return s.data
}

// GenerateID generates a unique ID
func GenerateID() string {
	return uuid.New().String()[:8]
}

// ==================== Expense Operations ====================

// AddExpense adds a new expense
func (s *Storage) AddExpense(amount float64, description string, category models.ExpenseCategory, date time.Time) (*models.Expense, error) {
	expense := models.Expense{
		ID:          GenerateID(),
		Amount:      amount,
		Description: description,
		Category:    category,
		Date:        date,
		CreatedAt:   time.Now(),
	}
	s.data.Expenses = append(s.data.Expenses, expense)
	return &expense, s.Save()
}

// GetExpenses returns all expenses
func (s *Storage) GetExpenses() []models.Expense {
	return s.data.Expenses
}

// DeleteExpense deletes an expense by ID
func (s *Storage) DeleteExpense(id string) error {
	for i, exp := range s.data.Expenses {
		if exp.ID == id {
			s.data.Expenses = append(s.data.Expenses[:i], s.data.Expenses[i+1:]...)
			return s.Save()
		}
	}
	return nil
}

// ==================== Debt Transaction Operations ====================

// AddDebtTransaction adds a new debt transaction
func (s *Storage) AddDebtTransaction(txType models.TransactionType, personName string, amount float64, description string, date time.Time, dueDate *time.Time) (*models.DebtTransaction, error) {
	tx := models.DebtTransaction{
		ID:          GenerateID(),
		Type:        txType,
		PersonName:  personName,
		Amount:      amount,
		Description: description,
		Date:        date,
		DueDate:     dueDate,
		IsSettled:   false,
		CreatedAt:   time.Now(),
	}
	s.data.DebtTransactions = append(s.data.DebtTransactions, tx)
	return &tx, s.Save()
}

// SettleDebtTransaction marks a transaction as settled
func (s *Storage) SettleDebtTransaction(id string) error {
	for i, tx := range s.data.DebtTransactions {
		if tx.ID == id {
			now := time.Now()
			s.data.DebtTransactions[i].IsSettled = true
			s.data.DebtTransactions[i].SettledDate = &now
			return s.Save()
		}
	}
	return nil
}

// PartialSettleDebt settles a specific amount for a person
// It settles transactions in order until the amount is covered
// Returns the actual amount settled
func (s *Storage) PartialSettleDebt(personName string, amount float64, settleType models.TransactionType) (float64, error) {
	var settled float64
	now := time.Now()

	for i, tx := range s.data.DebtTransactions {
		if tx.PersonName == personName && tx.Type == settleType && !tx.IsSettled {
			if settled >= amount {
				break
			}

			remaining := amount - settled
			if tx.Amount <= remaining {
				// Fully settle this transaction
				s.data.DebtTransactions[i].IsSettled = true
				s.data.DebtTransactions[i].SettledDate = &now
				settled += tx.Amount
			} else {
				// Partially settle - reduce amount
				s.data.DebtTransactions[i].Amount -= remaining
				settled += remaining
			}
		}
	}

	if settled > 0 {
		return settled, s.Save()
	}
	return 0, nil
}

// SettleAmountForPerson settles a specific amount for a person (handles both lent and borrowed)
// It calculates net balance and settles appropriately
// settlementNote is an optional description of why/how the settlement occurred
func (s *Storage) SettleAmountForPerson(personName string, amount float64, settlementNote string) (float64, error) {
	// Calculate current balances
	var totalLent, totalBorrowed float64
	for _, tx := range s.data.DebtTransactions {
		if tx.PersonName == personName && !tx.IsSettled {
			if tx.Type == models.Lent {
				totalLent += tx.Amount
			} else {
				totalBorrowed += tx.Amount
			}
		}
	}

	netBalance := totalLent - totalBorrowed
	now := time.Now()
	var settled float64

	if netBalance > 0 {
		// They owe us - settle from lent transactions
		remainingToSettle := amount
		if remainingToSettle == 0 || remainingToSettle > netBalance {
			remainingToSettle = netBalance
		}
		for i, tx := range s.data.DebtTransactions {
			if tx.PersonName == personName && tx.Type == models.Lent && !tx.IsSettled && remainingToSettle > 0 {
				if tx.Amount <= remainingToSettle {
					s.data.DebtTransactions[i].IsSettled = true
					s.data.DebtTransactions[i].SettledDate = &now
					s.data.DebtTransactions[i].SettlementAmount = tx.Amount
					s.data.DebtTransactions[i].SettlementNote = settlementNote
					settled += tx.Amount
					remainingToSettle -= tx.Amount
				} else {
					s.data.DebtTransactions[i].Amount -= remainingToSettle
					settled += remainingToSettle
					remainingToSettle = 0
				}
			}
		}
		// Also settle borrowed transactions up to the same amount (offsetting)
		offsetSettle := netBalance
		if offsetSettle == 0 {
			offsetSettle = netBalance
		}
		for i, tx := range s.data.DebtTransactions {
			if tx.PersonName == personName && tx.Type == models.Borrowed && !tx.IsSettled && offsetSettle > 0 {
				if tx.Amount <= offsetSettle {
					s.data.DebtTransactions[i].IsSettled = true
					s.data.DebtTransactions[i].SettledDate = &now
					s.data.DebtTransactions[i].SettlementAmount = tx.Amount
					s.data.DebtTransactions[i].SettlementNote = settlementNote
					offsetSettle -= tx.Amount
				} else {
					s.data.DebtTransactions[i].Amount -= offsetSettle
					offsetSettle = 0
				}
			}
		}
	} else if netBalance < 0 {
		// We owe them - settle from borrowed transactions
		remainingToSettle := amount
		if remainingToSettle == 0 || remainingToSettle > -netBalance {
			remainingToSettle = -netBalance
		}
		for i, tx := range s.data.DebtTransactions {
			if tx.PersonName == personName && tx.Type == models.Borrowed && !tx.IsSettled && remainingToSettle > 0 {
				if tx.Amount <= remainingToSettle {
					s.data.DebtTransactions[i].IsSettled = true
					s.data.DebtTransactions[i].SettledDate = &now
					s.data.DebtTransactions[i].SettlementAmount = tx.Amount
					s.data.DebtTransactions[i].SettlementNote = settlementNote
					settled += tx.Amount
					remainingToSettle -= tx.Amount
				} else {
					s.data.DebtTransactions[i].Amount -= remainingToSettle
					settled += remainingToSettle
					remainingToSettle = 0
				}
			}
		}
		// Also settle lent transactions up to the same amount (offsetting)
		offsetSettle := netBalance
		if offsetSettle == 0 {
			offsetSettle = -netBalance
		}
		for i, tx := range s.data.DebtTransactions {
			if tx.PersonName == personName && tx.Type == models.Lent && !tx.IsSettled && offsetSettle > 0 {
				if tx.Amount <= offsetSettle {
					s.data.DebtTransactions[i].IsSettled = true
					s.data.DebtTransactions[i].SettledDate = &now
					s.data.DebtTransactions[i].SettlementAmount = tx.Amount
					s.data.DebtTransactions[i].SettlementNote = settlementNote
					offsetSettle -= tx.Amount
				} else {
					s.data.DebtTransactions[i].Amount -= offsetSettle
					offsetSettle = 0
				}
			}
		}
	} else if netBalance == 0 {
		// Net is 0 but there might be unsettled transactions - settle all
		var hasUnsettled bool
		for i, tx := range s.data.DebtTransactions {
			if tx.PersonName == personName && !tx.IsSettled {
				s.data.DebtTransactions[i].IsSettled = true
				s.data.DebtTransactions[i].SettledDate = &now
				s.data.DebtTransactions[i].SettlementAmount = tx.Amount
				s.data.DebtTransactions[i].SettlementNote = settlementNote
				settled += tx.Amount
				hasUnsettled = true
			}
		}
		if !hasUnsettled {
			return 0, nil
		}
	}

	if settled > 0 {
		return settled, s.Save()
	}
	return 0, nil
}

// GetPersonNetBalance returns the net balance for a person
func (s *Storage) GetPersonNetBalance(personName string) float64 {
	var totalLent, totalBorrowed float64
	for _, tx := range s.data.DebtTransactions {
		if tx.PersonName == personName && !tx.IsSettled {
			if tx.Type == models.Lent {
				totalLent += tx.Amount
			} else {
				totalBorrowed += tx.Amount
			}
		}
	}
	return totalLent - totalBorrowed
}

// GetDebtTransactions returns all debt transactions
func (s *Storage) GetDebtTransactions() []models.DebtTransaction {
	return s.data.DebtTransactions
}

// GetUnsettledDebts returns unsettled debt transactions
func (s *Storage) GetUnsettledDebts() []models.DebtTransaction {
	var unsettled []models.DebtTransaction
	for _, tx := range s.data.DebtTransactions {
		if !tx.IsSettled {
			unsettled = append(unsettled, tx)
		}
	}
	return unsettled
}

// GetSettledDebts returns settled debt transactions
func (s *Storage) GetSettledDebts() []models.DebtTransaction {
	var settled []models.DebtTransaction
	for _, tx := range s.data.DebtTransactions {
		if tx.IsSettled {
			settled = append(settled, tx)
		}
	}
	return settled
}

// GetUnsettledDebtsForPerson returns unsettled debts for a specific person
func (s *Storage) GetUnsettledDebtsForPerson(personName string) []models.DebtTransaction {
	var debts []models.DebtTransaction
	for _, tx := range s.data.DebtTransactions {
		if tx.PersonName == personName && !tx.IsSettled {
			debts = append(debts, tx)
		}
	}
	return debts
}

// GetAllDebtsForPerson returns all debts (settled and unsettled) for a specific person
func (s *Storage) GetAllDebtsForPerson(personName string) []models.DebtTransaction {
	var debts []models.DebtTransaction
	for _, tx := range s.data.DebtTransactions {
		if tx.PersonName == personName {
			debts = append(debts, tx)
		}
	}
	return debts
}

// SettleTransaction settles a specific transaction by ID with a specific amount
// If amount is less than the full transaction amount, the transaction is split:
// - Original transaction remains with reduced amount
// - New settled transaction is created for the settled portion
func (s *Storage) SettleTransaction(id string, amount float64, note string) error {
	for i, tx := range s.data.DebtTransactions {
		if tx.ID == id {
			now := time.Now()

			if amount >= tx.Amount {
				// Full settlement - mark original as settled
				s.data.DebtTransactions[i].IsSettled = true
				s.data.DebtTransactions[i].SettledDate = &now
				s.data.DebtTransactions[i].SettlementAmount = tx.Amount
				s.data.DebtTransactions[i].SettlementNote = note
			} else {
				// Partial settlement - split the transaction
				// 1. Reduce original transaction amount
				s.data.DebtTransactions[i].Amount -= amount

				// 2. Create a new settled transaction for the settled portion
				settledTx := models.DebtTransaction{
					ID:               GenerateID(),
					Type:             tx.Type,
					PersonName:       tx.PersonName,
					Amount:           amount,
					Description:      tx.Description + " (partial settlement)",
					Date:             tx.Date,
					DueDate:          tx.DueDate,
					IsSettled:        true,
					SettledDate:      &now,
					SettlementAmount: amount,
					SettlementNote:   note,
					CreatedAt:        now,
				}
				s.data.DebtTransactions = append(s.data.DebtTransactions, settledTx)
			}

			return s.Save()
		}
	}
	return fmt.Errorf("transaction not found")
}

// ==================== Investment Operations ====================

// AddInvestment adds a new investment
func (s *Storage) AddInvestment(invType models.InvestmentType, name string, investedAmount, currentValue, units float64, purchaseDate time.Time, notes string) (*models.Investment, error) {
	inv := models.Investment{
		ID:             GenerateID(),
		Type:           invType,
		Name:           name,
		InvestedAmount: investedAmount,
		CurrentValue:   currentValue,
		Units:          units,
		PurchaseDate:   purchaseDate,
		Notes:          notes,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	s.data.Investments = append(s.data.Investments, inv)
	return &inv, s.Save()
}

// UpdateInvestmentValue updates the current value of an investment
func (s *Storage) UpdateInvestmentValue(id string, currentValue float64) error {
	for i, inv := range s.data.Investments {
		if inv.ID == id {
			s.data.Investments[i].CurrentValue = currentValue
			s.data.Investments[i].UpdatedAt = time.Now()
			return s.Save()
		}
	}
	return nil
}

// UpdateInvestment updates both invested amount and current value of an investment
func (s *Storage) UpdateInvestment(id string, investedAmount, currentValue float64) error {
	for i, inv := range s.data.Investments {
		if inv.ID == id {
			s.data.Investments[i].InvestedAmount = investedAmount
			s.data.Investments[i].CurrentValue = currentValue
			s.data.Investments[i].UpdatedAt = time.Now()
			return s.Save()
		}
	}
	return nil
}

// GetInvestments returns all investments
func (s *Storage) GetInvestments() []models.Investment {
	return s.data.Investments
}

// GetInvestmentsByType returns investments of a specific type
func (s *Storage) GetInvestmentsByType(invType models.InvestmentType) []models.Investment {
	var investments []models.Investment
	for _, inv := range s.data.Investments {
		if inv.Type == invType {
			investments = append(investments, inv)
		}
	}
	return investments
}

// DeleteInvestment deletes an investment by ID
func (s *Storage) DeleteInvestment(id string) error {
	for i, inv := range s.data.Investments {
		if inv.ID == id {
			s.data.Investments = append(s.data.Investments[:i], s.data.Investments[i+1:]...)
			return s.Save()
		}
	}
	return nil
}

// ==================== Savings Target Operations ====================

// AddSavingsTarget adds a new savings target
func (s *Storage) AddSavingsTarget(productName string, targetAmount float64, targetDate time.Time, description string) (*models.SavingsTarget, error) {
	target := models.SavingsTarget{
		ID:            GenerateID(),
		ProductName:   productName,
		TargetAmount:  targetAmount,
		CurrentAmount: 0,
		TargetDate:    targetDate,
		Description:   description,
		IsCompleted:   false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	s.data.SavingsTargets = append(s.data.SavingsTargets, target)
	return &target, s.Save()
}

// AddSavingsContribution adds a contribution to a savings target
func (s *Storage) AddSavingsContribution(targetID string, amount float64, notes string) (*models.SavingsContribution, error) {
	// Find and update the target
	var targetFound bool
	for i, target := range s.data.SavingsTargets {
		if target.ID == targetID {
			s.data.SavingsTargets[i].CurrentAmount += amount
			s.data.SavingsTargets[i].UpdatedAt = time.Now()
			if s.data.SavingsTargets[i].CurrentAmount >= s.data.SavingsTargets[i].TargetAmount {
				s.data.SavingsTargets[i].IsCompleted = true
			}
			targetFound = true
			break
		}
	}

	if !targetFound {
		return nil, nil
	}

	contribution := models.SavingsContribution{
		ID:        GenerateID(),
		TargetID:  targetID,
		Amount:    amount,
		Date:      time.Now(),
		Notes:     notes,
		CreatedAt: time.Now(),
	}
	s.data.SavingsContributions = append(s.data.SavingsContributions, contribution)
	return &contribution, s.Save()
}

// GetSavingsTargets returns all savings targets
func (s *Storage) GetSavingsTargets() []models.SavingsTarget {
	return s.data.SavingsTargets
}

// GetActiveSavingsTargets returns non-completed savings targets
func (s *Storage) GetActiveSavingsTargets() []models.SavingsTarget {
	var active []models.SavingsTarget
	for _, target := range s.data.SavingsTargets {
		if !target.IsCompleted {
			active = append(active, target)
		}
	}
	return active
}

// GetSavingsContributions returns contributions for a target
func (s *Storage) GetSavingsContributions(targetID string) []models.SavingsContribution {
	var contributions []models.SavingsContribution
	for _, c := range s.data.SavingsContributions {
		if c.TargetID == targetID {
			contributions = append(contributions, c)
		}
	}
	return contributions
}

// DeleteSavingsTarget deletes a savings target by ID
func (s *Storage) DeleteSavingsTarget(id string) error {
	for i, target := range s.data.SavingsTargets {
		if target.ID == id {
			s.data.SavingsTargets = append(s.data.SavingsTargets[:i], s.data.SavingsTargets[i+1:]...)
			return s.Save()
		}
	}
	return nil
}
