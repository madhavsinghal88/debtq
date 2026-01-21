package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/debtq/debtq/internal/config"
	"github.com/debtq/debtq/internal/models"
)

// ObsidianWriter handles writing markdown files to Obsidian vault
type ObsidianWriter struct {
	config *config.Config
}

// NewObsidianWriter creates a new ObsidianWriter
func NewObsidianWriter(cfg *config.Config) *ObsidianWriter {
	return &ObsidianWriter{config: cfg}
}

// EnsureDirs ensures all required directories exist
func (o *ObsidianWriter) EnsureDirs() error {
	if err := os.MkdirAll(o.config.ObsidianVaultPath, 0755); err != nil {
		return err
	}
	return nil
}

// SyncAllNotes syncs all data to Obsidian vault as summarized files
func (o *ObsidianWriter) SyncAllNotes(data *models.Data) error {
	if err := o.EnsureDirs(); err != nil {
		return err
	}

	// Write main dashboard
	if err := o.writeDashboard(data); err != nil {
		return err
	}

	// Write expenses summary
	if err := o.writeExpensesSummary(data); err != nil {
		return err
	}

	// Write debts summary (grouped by person)
	if err := o.writeDebtsSummary(data); err != nil {
		return err
	}

	// Write net worth summary
	if err := o.writeNetWorthSummary(data); err != nil {
		return err
	}

	// Write savings summary
	if err := o.writeSavingsSummary(data); err != nil {
		return err
	}

	return nil
}

// writeDashboard writes the main dashboard file
func (o *ObsidianWriter) writeDashboard(data *models.Data) error {
	now := time.Now()

	type Dashboard struct {
		NetWorth           float64
		TotalBorrowed      float64
		TotalLent          float64
		NetDebtPosition    float64
		MonthlyExpenses    float64
		TotalExpenses      float64
		ActiveSavingsGoals int
		TotalSavingsTarget float64
		TotalSaved         float64
		SavingsProgress    float64
		UpdatedAt          time.Time
	}

	var totalExpenses, totalSavingsTarget, totalSaved float64
	var activeSavings int

	for _, e := range data.Expenses {
		totalExpenses += e.Amount
	}

	for _, t := range data.SavingsTargets {
		if !t.IsCompleted {
			activeSavings++
		}
		totalSavingsTarget += t.TargetAmount
		totalSaved += t.CurrentAmount
	}

	savingsProgress := float64(0)
	if totalSavingsTarget > 0 {
		savingsProgress = (totalSaved / totalSavingsTarget) * 100
	}

	dashboard := Dashboard{
		NetWorth:           data.NetWorth(),
		TotalBorrowed:      data.TotalBorrowed(),
		TotalLent:          data.TotalLent(),
		NetDebtPosition:    data.TotalLent() - data.TotalBorrowed(),
		MonthlyExpenses:    data.MonthlyExpenses(now.Year(), now.Month()),
		TotalExpenses:      totalExpenses,
		ActiveSavingsGoals: activeSavings,
		TotalSavingsTarget: totalSavingsTarget,
		TotalSaved:         totalSaved,
		SavingsProgress:    savingsProgress,
		UpdatedAt:          now,
	}

	tmpl := `---
tags: [debtq, dashboard, finance]
updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}
---

# DebtQ - Financial Dashboard

> Last Updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}

## Quick Overview

| Category | Amount |
|----------|--------|
| **Net Worth** | {{printf "%.2f" .NetWorth}} |
| **Net Debt Position** | {{printf "%.2f" .NetDebtPosition}} |
| **This Month Expenses** | {{printf "%.2f" .MonthlyExpenses}} |

---

## Net Worth
Total investments value: **{{printf "%.2f" .NetWorth}}**

[[NetWorth|View Details →]]

---

## Debts & Lending

| Metric | Amount |
|--------|--------|
| Total Lent (others owe you) | {{printf "%.2f" .TotalLent}} |
| Total Borrowed (you owe) | {{printf "%.2f" .TotalBorrowed}} |
| **Net Position** | {{printf "%.2f" .NetDebtPosition}} |

[[Debts|View Details →]]

---

## Expenses

| Metric | Amount |
|--------|--------|
| This Month | {{printf "%.2f" .MonthlyExpenses}} |
| All Time Total | {{printf "%.2f" .TotalExpenses}} |

[[Expenses|View Details →]]

---

## Savings Goals

| Metric | Value |
|--------|-------|
| Active Goals | {{.ActiveSavingsGoals}} |
| Total Target | {{printf "%.2f" .TotalSavingsTarget}} |
| Total Saved | {{printf "%.2f" .TotalSaved}} |
| Progress | {{printf "%.1f" .SavingsProgress}}% |

[[Savings|View Details →]]
`

	return o.writeNoteWithFuncs("", "Dashboard.md", tmpl, dashboard)
}

// writeExpensesSummary writes expenses grouped by month and category
func (o *ObsidianWriter) writeExpensesSummary(data *models.Data) error {
	type MonthData struct {
		Month      string
		Total      float64
		ByCategory map[string]float64
		Expenses   []models.Expense
	}

	type ExpensesSummary struct {
		Months     []MonthData
		TotalAll   float64
		ByCategory map[string]float64
		UpdatedAt  time.Time
	}

	// Group expenses by month
	monthMap := make(map[string]*MonthData)
	var monthOrder []string
	totalByCategory := make(map[string]float64)
	var totalAll float64

	for _, exp := range data.Expenses {
		monthKey := exp.Date.Format("2006-01")
		if _, exists := monthMap[monthKey]; !exists {
			monthMap[monthKey] = &MonthData{
				Month:      exp.Date.Format("January 2006"),
				Total:      0,
				ByCategory: make(map[string]float64),
				Expenses:   []models.Expense{},
			}
			monthOrder = append(monthOrder, monthKey)
		}
		monthMap[monthKey].Total += exp.Amount
		monthMap[monthKey].ByCategory[string(exp.Category)] += exp.Amount
		monthMap[monthKey].Expenses = append(monthMap[monthKey].Expenses, exp)
		totalByCategory[string(exp.Category)] += exp.Amount
		totalAll += exp.Amount
	}

	// Sort months in reverse order (newest first)
	sort.Sort(sort.Reverse(sort.StringSlice(monthOrder)))

	var months []MonthData
	for _, key := range monthOrder {
		months = append(months, *monthMap[key])
	}

	summary := ExpensesSummary{
		Months:     months,
		TotalAll:   totalAll,
		ByCategory: totalByCategory,
		UpdatedAt:  time.Now(),
	}

	tmpl := `---
tags: [debtq, expenses, finance]
updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}
---

# Expenses Summary

> Last Updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}

## Total: {{printf "%.2f" .TotalAll}}

### By Category (All Time)

| Category | Amount |
|----------|--------|
{{- range $cat, $amt := .ByCategory}}
| {{$cat}} | {{printf "%.2f" $amt}} |
{{- end}}

---
{{range .Months}}
## {{.Month}}

**Total: {{printf "%.2f" .Total}}**

| Date | Description | Category | Amount |
|------|-------------|----------|--------|
{{- range .Expenses}}
| {{.Date.Format "02"}} | {{.Description}} | {{.Category}} | {{printf "%.2f" .Amount}} |
{{- end}}

{{end}}
`

	return o.writeNoteWithFuncs("", "Expenses.md", tmpl, summary)
}

// writeDebtsSummary writes debts grouped by person
func (o *ObsidianWriter) writeDebtsSummary(data *models.Data) error {
	type PersonDebt struct {
		Name          string
		TotalLent     float64
		TotalBorrowed float64
		NetBalance    float64
		LentTxns      []models.DebtTransaction
		BorrowedTxns  []models.DebtTransaction
	}

	type DebtsSummary struct {
		People        []PersonDebt
		TotalLent     float64
		TotalBorrowed float64
		NetPosition   float64
		UpdatedAt     time.Time
	}

	// Group by person
	personMap := make(map[string]*PersonDebt)
	var personOrder []string

	for _, tx := range data.DebtTransactions {
		if tx.IsSettled {
			continue
		}
		if _, exists := personMap[tx.PersonName]; !exists {
			personMap[tx.PersonName] = &PersonDebt{
				Name:          tx.PersonName,
				TotalLent:     0,
				TotalBorrowed: 0,
				LentTxns:      []models.DebtTransaction{},
				BorrowedTxns:  []models.DebtTransaction{},
			}
			personOrder = append(personOrder, tx.PersonName)
		}
		if tx.Type == models.Lent {
			personMap[tx.PersonName].TotalLent += tx.Amount
			personMap[tx.PersonName].LentTxns = append(personMap[tx.PersonName].LentTxns, tx)
		} else {
			personMap[tx.PersonName].TotalBorrowed += tx.Amount
			personMap[tx.PersonName].BorrowedTxns = append(personMap[tx.PersonName].BorrowedTxns, tx)
		}
	}

	var people []PersonDebt
	for _, name := range personOrder {
		p := personMap[name]
		p.NetBalance = p.TotalLent - p.TotalBorrowed
		people = append(people, *p)
	}

	// Sort by net balance (highest owed to you first)
	sort.Slice(people, func(i, j int) bool {
		return people[i].NetBalance > people[j].NetBalance
	})

	summary := DebtsSummary{
		People:        people,
		TotalLent:     data.TotalLent(),
		TotalBorrowed: data.TotalBorrowed(),
		NetPosition:   data.TotalLent() - data.TotalBorrowed(),
		UpdatedAt:     time.Now(),
	}

	tmpl := `---
tags: [debtq, debts, lending, finance]
updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}
---

# Debts & Lending Summary

> Last Updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}

## Overview

| Metric | Amount |
|--------|--------|
| Total Lent (others owe you) | {{printf "%.2f" .TotalLent}} |
| Total Borrowed (you owe) | {{printf "%.2f" .TotalBorrowed}} |
| **Net Position** | {{printf "%.2f" .NetPosition}} |

---

## By Person
{{if not .People}}
*No pending debts*
{{end}}
{{range .People}}
### {{.Name}}

{{if gt .NetBalance 0.0}}**Owes you: {{printf "%.2f" .NetBalance}}**{{else if lt .NetBalance 0.0}}**You owe: {{printf "%.2f" (neg .NetBalance)}}**{{else}}**Settled**{{end}}

{{if .LentTxns}}
**Lent:**
| Date | Amount | Reason |
|------|--------|--------|
{{- range .LentTxns}}
| {{.Date.Format "2006-01-02"}} | +{{printf "%.2f" .Amount}} | {{.Description}} |
{{- end}}
{{end}}
{{if .BorrowedTxns}}
**Borrowed:**
| Date | Amount | Reason |
|------|--------|--------|
{{- range .BorrowedTxns}}
| {{.Date.Format "2006-01-02"}} | -{{printf "%.2f" .Amount}} | {{.Description}} |
{{- end}}
{{end}}

---
{{end}}
`

	return o.writeNoteWithFuncs("", "Debts.md", tmpl, summary)
}

// writeNetWorthSummary writes investments summary
func (o *ObsidianWriter) writeNetWorthSummary(data *models.Data) error {
	type InvestmentGroup struct {
		Type           string
		TotalInvested  float64
		TotalCurrent   float64
		Gain           float64
		GainPercentage float64
		Investments    []models.Investment
	}

	type NetWorthSummary struct {
		Groups         []InvestmentGroup
		TotalInvested  float64
		TotalCurrent   float64
		TotalGain      float64
		GainPercentage float64
		UpdatedAt      time.Time
	}

	// Group by investment type
	typeMap := make(map[string]*InvestmentGroup)
	var typeOrder []string
	var totalInvested, totalCurrent float64

	for _, inv := range data.Investments {
		typeKey := string(inv.Type)
		if _, exists := typeMap[typeKey]; !exists {
			typeMap[typeKey] = &InvestmentGroup{
				Type:        typeKey,
				Investments: []models.Investment{},
			}
			typeOrder = append(typeOrder, typeKey)
		}
		typeMap[typeKey].TotalInvested += inv.InvestedAmount
		typeMap[typeKey].TotalCurrent += inv.CurrentValue
		typeMap[typeKey].Investments = append(typeMap[typeKey].Investments, inv)
		totalInvested += inv.InvestedAmount
		totalCurrent += inv.CurrentValue
	}

	var groups []InvestmentGroup
	for _, key := range typeOrder {
		g := typeMap[key]
		g.Gain = g.TotalCurrent - g.TotalInvested
		if g.TotalInvested > 0 {
			g.GainPercentage = (g.Gain / g.TotalInvested) * 100
		}
		groups = append(groups, *g)
	}

	// Sort by current value (highest first)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].TotalCurrent > groups[j].TotalCurrent
	})

	totalGain := totalCurrent - totalInvested
	gainPercentage := float64(0)
	if totalInvested > 0 {
		gainPercentage = (totalGain / totalInvested) * 100
	}

	summary := NetWorthSummary{
		Groups:         groups,
		TotalInvested:  totalInvested,
		TotalCurrent:   totalCurrent,
		TotalGain:      totalGain,
		GainPercentage: gainPercentage,
		UpdatedAt:      time.Now(),
	}

	tmpl := `---
tags: [debtq, networth, investments, finance]
updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}
---

# My Net Worth

> Last Updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}

## Overview

| Metric | Value |
|--------|-------|
| Total Invested | {{printf "%.2f" .TotalInvested}} |
| Current Value | {{printf "%.2f" .TotalCurrent}} |
| Total Gain/Loss | {{printf "%.2f" .TotalGain}} |
| Return | {{printf "%.2f" .GainPercentage}}% |

---

## By Investment Type

| Type | Invested | Current | Gain/Loss | Return % |
|------|----------|---------|-----------|----------|
{{- range .Groups}}
| **{{.Type}}** | {{printf "%.2f" .TotalInvested}} | {{printf "%.2f" .TotalCurrent}} | {{printf "%.2f" .Gain}} | {{printf "%.2f" .GainPercentage}}% |
{{- end}}

---
{{range .Groups}}
## {{.Type}}

| Name | Invested | Current | Gain/Loss | Return % |
|------|----------|---------|-----------|----------|
{{- range .Investments}}
| {{.Name}} | {{printf "%.2f" .InvestedAmount}} | {{printf "%.2f" .CurrentValue}} | {{printf "%.2f" (sub .CurrentValue .InvestedAmount)}} | {{if gt .InvestedAmount 0}}{{printf "%.2f" (percentage .CurrentValue .InvestedAmount)}}%{{else}}N/A{{end}} |
{{- end}}

{{end}}
`

	return o.writeNoteWithFuncs("", "NetWorth.md", tmpl, summary)
}

// writeSavingsSummary writes savings goals summary
func (o *ObsidianWriter) writeSavingsSummary(data *models.Data) error {
	type SavingsSummary struct {
		ActiveGoals    []models.SavingsTarget
		CompletedGoals []models.SavingsTarget
		TotalTarget    float64
		TotalSaved     float64
		Progress       float64
		UpdatedAt      time.Time
	}

	var active, completed []models.SavingsTarget
	var totalTarget, totalSaved float64

	for _, t := range data.SavingsTargets {
		totalTarget += t.TargetAmount
		totalSaved += t.CurrentAmount
		if t.IsCompleted {
			completed = append(completed, t)
		} else {
			active = append(active, t)
		}
	}

	progress := float64(0)
	if totalTarget > 0 {
		progress = (totalSaved / totalTarget) * 100
	}

	summary := SavingsSummary{
		ActiveGoals:    active,
		CompletedGoals: completed,
		TotalTarget:    totalTarget,
		TotalSaved:     totalSaved,
		Progress:       progress,
		UpdatedAt:      time.Now(),
	}

	tmpl := `---
tags: [debtq, savings, goals, finance]
updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}
---

# Savings Goals

> Last Updated: {{.UpdatedAt.Format "2006-01-02 15:04:05"}}

## Overview

| Metric | Value |
|--------|-------|
| Total Target | {{printf "%.2f" .TotalTarget}} |
| Total Saved | {{printf "%.2f" .TotalSaved}} |
| Overall Progress | {{printf "%.1f" .Progress}}% |
| Active Goals | {{len .ActiveGoals}} |
| Completed Goals | {{len .CompletedGoals}} |

---

## Active Goals
{{if not .ActiveGoals}}
*No active savings goals*
{{end}}
{{range .ActiveGoals}}
### {{.ProductName}}

| Metric | Value |
|--------|-------|
| Target | {{printf "%.2f" .TargetAmount}} |
| Saved | {{printf "%.2f" .CurrentAmount}} |
| Remaining | {{printf "%.2f" (sub .TargetAmount .CurrentAmount)}} |
| Progress | {{printf "%.1f" (percentage .CurrentAmount .TargetAmount)}}% |
| Target Date | {{.TargetDate.Format "2006-01-02"}} |
| Days Left | {{daysRemaining .TargetDate}} |
| Monthly Required | {{printf "%.2f" (monthlyRequired .TargetAmount .CurrentAmount .TargetDate)}} |

` + "```" + `
{{progressBar .CurrentAmount .TargetAmount 30}}
` + "```" + `

{{if .Description}}*{{.Description}}*{{end}}

---
{{end}}

{{if .CompletedGoals}}
## Completed Goals

| Product | Target | Saved | Completed |
|---------|--------|-------|-----------|
{{- range .CompletedGoals}}
| {{.ProductName}} | {{printf "%.2f" .TargetAmount}} | {{printf "%.2f" .CurrentAmount}} | ✅ |
{{- end}}
{{end}}
`

	return o.writeNoteWithFuncs("", "Savings.md", tmpl, summary)
}

// Helper functions

func (o *ObsidianWriter) writeNoteWithFuncs(subdir, filename, tmplStr string, data interface{}) error {
	funcMap := template.FuncMap{
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"neg": func(a float64) float64 {
			return -a
		},
		"gt": func(a, b float64) bool {
			return a > b
		},
		"lt": func(a, b float64) bool {
			return a < b
		},
		"percentage": func(current, total float64) float64 {
			if total == 0 {
				return 0
			}
			return ((current-total)/total)*100 + 100
		},
		"daysRemaining": func(targetDate time.Time) int {
			days := int(time.Until(targetDate).Hours() / 24)
			if days < 0 {
				return 0
			}
			return days
		},
		"monthlyRequired": func(target, current float64, targetDate time.Time) float64 {
			remaining := target - current
			if remaining <= 0 {
				return 0
			}
			months := float64(time.Until(targetDate).Hours()) / (24 * 30)
			if months <= 0 {
				return remaining
			}
			return remaining / months
		},
		"progressBar": func(current, total float64, width int) string {
			if total == 0 {
				return ""
			}
			pct := current / total
			if pct > 1 {
				pct = 1
			}
			filled := int(pct * float64(width))
			empty := width - filled
			bar := ""
			for i := 0; i < filled; i++ {
				bar += "█"
			}
			for i := 0; i < empty; i++ {
				bar += "░"
			}
			return fmt.Sprintf("%s %.1f%%", bar, pct*100)
		},
	}

	tmpl, err := template.New("note").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	var filePath string
	if subdir == "" {
		filePath = filepath.Join(o.config.ObsidianVaultPath, filename)
	} else {
		filePath = filepath.Join(o.config.ObsidianVaultPath, subdir, filename)
	}
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func sanitizeFilename(s string) string {
	result := ""
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result += string(c)
		} else if c == ' ' {
			result += "-"
		}
	}
	return result
}
