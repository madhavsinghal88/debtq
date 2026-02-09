package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/debtq/debtq/internal/config"
	"github.com/debtq/debtq/internal/models"
	"github.com/debtq/debtq/internal/storage"
)

// View represents different screens
type View int

const (
	ViewMain View = iota
	ViewExpenses
	ViewAddExpense
	ViewDebts
	ViewAddDebt
	ViewSelectTransactionToSettle
	ViewSettleDebt
	ViewDebtHistory
	ViewPersonHistory
	ViewNetWorth
	ViewAddInvestment
	ViewUpdateInvestment
	ViewConfirmDelete
	ViewSavings
	ViewAddSavingsTarget
	ViewAddContribution
	ViewStats
	ViewSettings
)

// Model is the main application model
type Model struct {
	config         *config.Config
	storage        *storage.Storage
	obsidian       *storage.ObsidianWriter
	currentView    View
	previousView   View
	cursor         int
	inputs         []textinput.Model
	focusIndex     int
	message        string
	messageType    string // "success", "error", "info"
	selectedID     string
	selectedPerson string
	width          int
	height         int
}

// New creates a new TUI model
func New(cfg *config.Config, store *storage.Storage) *Model {
	return &Model{
		config:      cfg,
		storage:     store,
		obsidian:    storage.NewObsidianWriter(cfg),
		currentView: ViewMain,
		cursor:      0,
		width:       80,
		height:      24,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Clear message on any key press
		if m.message != "" && msg.String() != "enter" {
			m.message = ""
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.currentView == ViewMain {
				return m, tea.Quit
			}
			// For views that have their own back navigation, let them handle it
			if m.currentView == ViewPersonHistory || m.currentView == ViewDebtHistory {
				// Let view-specific handler deal with it
				break
			}
			m.currentView = ViewMain
			m.cursor = 0
			m.inputs = nil
			return m, nil

		case "esc":
			if m.currentView == ViewMain {
				return m, nil
			}
			// For views that have their own back navigation, let them handle it
			if m.currentView == ViewPersonHistory || m.currentView == ViewDebtHistory {
				// Let view-specific handler deal with it
				break
			}
			m.currentView = ViewMain
			m.cursor = 0
			m.inputs = nil
			return m, nil
		}

		// Handle view-specific updates
		switch m.currentView {
		case ViewMain:
			return m.updateMainView(msg)
		case ViewExpenses:
			return m.updateExpensesView(msg)
		case ViewAddExpense:
			return m.updateAddExpenseView(msg)
		case ViewDebts:
			return m.updateDebtsView(msg)
		case ViewAddDebt:
			return m.updateAddDebtView(msg)
		case ViewSelectTransactionToSettle:
			return m.updateSelectTransactionToSettleView(msg)
		case ViewSettleDebt:
			return m.updateSettleDebtView(msg)
		case ViewDebtHistory:
			return m.updateDebtHistoryView(msg)
		case ViewPersonHistory:
			return m.updatePersonHistoryView(msg)
		case ViewNetWorth:
			return m.updateNetWorthView(msg)
		case ViewAddInvestment:
			return m.updateAddInvestmentView(msg)
		case ViewUpdateInvestment:
			return m.updateUpdateInvestmentView(msg)
		case ViewConfirmDelete:
			return m.updateConfirmDeleteView(msg)
		case ViewSavings:
			return m.updateSavingsView(msg)
		case ViewAddSavingsTarget:
			return m.updateAddSavingsTargetView(msg)
		case ViewAddContribution:
			return m.updateAddContributionView(msg)
		case ViewStats:
			return m.updateStatsView(msg)
		}
	}

	// Update text inputs if any
	if len(m.inputs) > 0 {
		cmds := make([]tea.Cmd, len(m.inputs))
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	var content string

	switch m.currentView {
	case ViewMain:
		content = m.viewMain()
	case ViewExpenses:
		content = m.viewExpenses()
	case ViewAddExpense:
		content = m.viewAddExpense()
	case ViewDebts:
		content = m.viewDebts()
	case ViewAddDebt:
		content = m.viewAddDebt()
	case ViewSelectTransactionToSettle:
		content = m.viewSelectTransactionToSettle()
	case ViewSettleDebt:
		content = m.viewSettleDebt()
	case ViewDebtHistory:
		content = m.viewDebtHistory()
	case ViewPersonHistory:
		content = m.viewPersonHistory()
	case ViewNetWorth:
		content = m.viewNetWorth()
	case ViewAddInvestment:
		content = m.viewAddInvestment()
	case ViewUpdateInvestment:
		content = m.viewUpdateInvestment()
	case ViewConfirmDelete:
		content = m.viewConfirmDelete()
	case ViewSavings:
		content = m.viewSavings()
	case ViewAddSavingsTarget:
		content = m.viewAddSavingsTarget()
	case ViewAddContribution:
		content = m.viewAddContribution()
	case ViewStats:
		content = m.viewStats()
	default:
		content = m.viewMain()
	}

	// Add message if present
	if m.message != "" {
		var msgStyle lipgloss.Style
		switch m.messageType {
		case "success":
			msgStyle = SuccessStyle
		case "error":
			msgStyle = ErrorStyle
		default:
			msgStyle = MutedStyle
		}
		content += "\n" + msgStyle.Render(m.message)
	}

	return content
}

// Main menu view - Beautiful centered layout with icons
func (m Model) viewMain() string {
	// Build the app header with branding
	header := AppTitleStyle.Width(70).Render("💰 DebtQ - Personal Money Tracker")
	subtitle := AppSubtitleStyle.Width(70).Render("Track expenses, debts, investments & savings goals")

	// Menu items with icons and descriptions
	type menuItem struct {
		icon        string
		title       string
		description string
	}

	menuItems := []menuItem{
		{"💳", "Expenses", "Track your daily spending"},
		{"🤝", "Borrowing & Lending", "Manage debts with friends"},
		{"📈", "My Net Worth", "Monitor investments & assets"},
		{"🎯", "Savings Goals", "Save for what matters"},
		{"📊", "Stats & Dashboard", "View financial insights"},
		{"📝", "Sync to Obsidian", "Export to markdown"},
		{"👋", "Quit", "Exit the application"},
	}

	// Build menu with cards
	var menuContent strings.Builder
	menuContent.WriteString("\n")

	for i, item := range menuItems {
		// Create card for each menu item
		var card lipgloss.Style
		if i == m.cursor {
			card = CardSelectedStyle.Width(66)
		} else {
			card = CardStyle.Width(66)
		}

		// Icon with color
		iconStyle := lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Width(3)

		// Title styling
		titleStyle := lipgloss.NewStyle().
			Foreground(Gray100).
			Bold(true).
			Width(25)

		// Description styling
		descStyle := lipgloss.NewStyle().
			Foreground(Gray500).
			Width(35)

		// Selection indicator
		indicator := "  "
		if i == m.cursor {
			indicator = "▸ "
			iconStyle = iconStyle.Foreground(Primary)
			titleStyle = titleStyle.Foreground(Primary)
		}

		// Build row
		row := fmt.Sprintf("%s%s %s %s",
			indicator,
			iconStyle.Render(item.icon),
			titleStyle.Render(item.title),
			descStyle.Render(item.description),
		)

		menuContent.WriteString(card.Render(row))
		menuContent.WriteString("\n")
	}

	// Help footer with styled keys
	helpBar := lipgloss.NewStyle().
		Background(Gray800).
		Foreground(Gray400).
		Padding(0, 2).
		Width(66).
		Render("↑/↓ Navigate • Enter Select • q Quit")

	// Quick stats preview
	data := m.storage.GetData()
	statsSection := ""
	if len(data.Expenses) > 0 || len(data.DebtTransactions) > 0 {
		statsContent := fmt.Sprintf("📊 %d Expenses  •  💸 %d Debts  •  📈 %d Investments",
			len(data.Expenses),
			len(data.DebtTransactions),
			len(data.Investments))
		statsSection = lipgloss.NewStyle().
			Foreground(Gray500).
			MarginTop(1).
			Render(statsContent) + "\n"
	}

	// Combine everything
	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		subtitle,
		menuContent.String(),
		statsSection,
		helpBar,
	)

	return content
}

func (m *Model) updateMainView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	menuLen := 7

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < menuLen-1 {
			m.cursor++
		}
	case "enter":
		switch m.cursor {
		case 0:
			m.currentView = ViewExpenses
			m.cursor = 0
		case 1:
			m.currentView = ViewDebts
			m.cursor = 0
		case 2:
			m.currentView = ViewNetWorth
			m.cursor = 0
		case 3:
			m.currentView = ViewSavings
			m.cursor = 0
		case 4:
			m.currentView = ViewStats
			m.cursor = 0
		case 5:
			// Sync to Obsidian
			if err := m.obsidian.SyncAllNotes(m.storage.GetData()); err != nil {
				m.message = "Error syncing: " + err.Error()
				m.messageType = "error"
			} else {
				m.message = "Successfully synced to Obsidian!"
				m.messageType = "success"
			}
		case 6:
			return m, tea.Quit
		}
	}

	return m, nil
}

// Expenses view
func (m Model) viewExpenses() string {
	title := TitleStyle.Render("  Expenses")

	expenses := m.storage.GetExpenses()

	var content string
	if len(expenses) == 0 {
		content = MutedStyle.Render("\n  No expenses recorded yet.\n")
	} else {
		content = "\n"
		// Show last 10 expenses
		start := 0
		if len(expenses) > 10 {
			start = len(expenses) - 10
		}
		for i := len(expenses) - 1; i >= start; i-- {
			exp := expenses[i]
			cursor := "  "
			if i-start == m.cursor {
				cursor = "▸ "
			}
			line := fmt.Sprintf("%s%s  %s  %s  %s",
				cursor,
				exp.Date.Format("2006-01-02"),
				TableCellStyle.Width(15).Render(truncate(exp.Description, 15)),
				TableCellStyle.Width(12).Render(string(exp.Category)),
				FormatAmount(exp.Amount, m.config.Currency),
			)
			content += line + "\n"
		}
	}

	// Calculate totals
	data := m.storage.GetData()
	now := time.Now()
	monthlyTotal := data.MonthlyExpenses(now.Year(), now.Month())

	stats := fmt.Sprintf("\n  This Month: %s", FormatAmountPlain(monthlyTotal, m.config.Currency))

	help := HelpStyle.Render("\n  a: Add expense • d: Delete • Enter: Details • Esc: Back")

	return BoxStyle.Render(title + content + stats + help)
}

func (m *Model) updateExpensesView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	expenses := m.storage.GetExpenses()
	maxCursor := len(expenses) - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "a":
		m.currentView = ViewAddExpense
		m.initExpenseInputs()
	case "d":
		if len(expenses) > 0 {
			idx := len(expenses) - 1 - m.cursor
			if idx >= 0 && idx < len(expenses) {
				m.storage.DeleteExpense(expenses[idx].ID)
				m.message = "Expense deleted"
				m.messageType = "success"
				if m.cursor > 0 {
					m.cursor--
				}
			}
		}
	}

	return m, nil
}

func (m *Model) initExpenseInputs() {
	m.inputs = make([]textinput.Model, 4)

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Amount"
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Description"

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "Category (food/transport/shopping/utilities/health/other)"

	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "Date (YYYY-MM-DD, leave empty for today)"

	m.focusIndex = 0
}

func (m Model) viewAddExpense() string {
	title := TitleStyle.Render("  Add Expense")

	var content string
	labels := []string{"Amount:", "Description:", "Category:", "Date:"}
	hints := []string{
		"",
		"",
		"Options: food, transport, shopping, utilities, health, entertainment, education, other",
		"Format: YYYY-MM-DD (leave empty for today)",
	}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		} else {
			content += MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		}
	}

	help := HelpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateAddExpenseView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		amount, err := strconv.ParseFloat(m.inputs[0].Value(), 64)
		if err != nil {
			m.message = "Invalid amount"
			m.messageType = "error"
			return m, nil
		}

		description := m.inputs[1].Value()
		if description == "" {
			m.message = "Description is required"
			m.messageType = "error"
			return m, nil
		}

		category := models.ExpenseCategory(m.inputs[2].Value())
		if category == "" {
			category = models.CategoryOther
		}

		date := time.Now()
		if m.inputs[3].Value() != "" {
			date, err = time.Parse("2006-01-02", m.inputs[3].Value())
			if err != nil {
				m.message = "Invalid date format"
				m.messageType = "error"
				return m, nil
			}
		}

		_, err = m.storage.AddExpense(amount, description, category, date)
		if err != nil {
			m.message = "Error saving expense: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Expense added successfully!"
		m.messageType = "success"
		m.currentView = ViewExpenses
		m.inputs = nil
		m.cursor = 0
		return m, nil
	case "+":
		if m.focusIndex == 0 && len(m.inputs) > 0 {
			currentValue := m.inputs[0].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[0].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	// Update text input
	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

// Debts view - Beautiful card-based layout
func (m Model) viewDebts() string {
	// Header with icon
	header := AppTitleStyle.Width(70).Render("🤝 Borrowing & Lending")

	debts := m.storage.GetUnsettledDebts()
	data := m.storage.GetData()

	var content strings.Builder

	if len(debts) == 0 {
		emptyState := CardStyle.Width(66).Render(
			lipgloss.NewStyle().
				Foreground(Gray500).
				Align(lipgloss.Center).
				Render("\n  📭 No unsettled debts\n  All caught up!\n"),
		)
		content.WriteString(emptyState)
	} else {
		// Group debts by person
		type personGroup struct {
			name          string
			totalLent     float64
			totalBorrowed float64
			lentDebts     []models.DebtTransaction
			borrowedDebts []models.DebtTransaction
		}

		groupMap := make(map[string]*personGroup)
		var groupOrder []string

		for _, debt := range debts {
			key := debt.PersonName
			if _, exists := groupMap[key]; !exists {
				groupMap[key] = &personGroup{
					name:          debt.PersonName,
					totalLent:     0,
					totalBorrowed: 0,
					lentDebts:     []models.DebtTransaction{},
					borrowedDebts: []models.DebtTransaction{},
				}
				groupOrder = append(groupOrder, key)
			}
			if debt.Type == models.Lent {
				groupMap[key].totalLent += debt.Amount
				groupMap[key].lentDebts = append(groupMap[key].lentDebts, debt)
			} else {
				groupMap[key].totalBorrowed += debt.Amount
				groupMap[key].borrowedDebts = append(groupMap[key].borrowedDebts, debt)
			}
		}

		content.WriteString("\n")
		visibleIndex := 0
		for _, key := range groupOrder {
			group := groupMap[key]
			netBalance := group.totalLent - group.totalBorrowed

			if netBalance == 0 {
				continue
			}

			// Create person card
			var card lipgloss.Style
			if visibleIndex == m.cursor {
				card = CardSelectedStyle.Width(66)
			} else {
				card = CardStyle.Width(66)
			}

			// Person header with status
			var statusIcon, statusText string
			if netBalance > 0 {
				statusIcon = "📥"
				statusText = RenderBadge(fmt.Sprintf("Owes you %s", FormatMoneyPlain(netBalance, m.config.Currency)), "success")
			} else {
				statusIcon = "📤"
				statusText = RenderBadge(fmt.Sprintf("You owe %s", FormatMoneyPlain(-netBalance, m.config.Currency)), "danger")
			}

			indicator := "  "
			if visibleIndex == m.cursor {
				indicator = "▸ "
			}

			personHeader := fmt.Sprintf("%s%s %s %s",
				indicator,
				statusIcon,
				lipgloss.NewStyle().Bold(true).Foreground(Gray100).Render(group.name),
				statusText,
			)

			var cardContent strings.Builder
			cardContent.WriteString(personHeader + "\n\n")

			// Show lent transactions
			if len(group.lentDebts) > 0 {
				cardContent.WriteString(lipgloss.NewStyle().Foreground(Secondary).Bold(true).Render("  📤 Lent") + "\n")
				for _, debt := range group.lentDebts {
					reason := debt.Description
					if reason == "" {
						reason = "(no description)"
					}
					line := fmt.Sprintf("    %s  %s  %s\n",
						MoneyPositiveStyle.Render("+"+FormatMoneyPlain(debt.Amount, m.config.Currency)),
						lipgloss.NewStyle().Foreground(Gray300).Render(truncate(reason, 25)),
						lipgloss.NewStyle().Foreground(Gray500).Render(debt.Date.Format("Jan 02")),
					)
					cardContent.WriteString(line)
				}
			}

			// Show borrowed transactions
			if len(group.borrowedDebts) > 0 {
				cardContent.WriteString(lipgloss.NewStyle().Foreground(Danger).Bold(true).Render("  📥 Borrowed") + "\n")
				for _, debt := range group.borrowedDebts {
					reason := debt.Description
					if reason == "" {
						reason = "(no description)"
					}
					line := fmt.Sprintf("    %s  %s  %s\n",
						MoneyNegativeStyle.Render("-"+FormatMoneyPlain(debt.Amount, m.config.Currency)),
						lipgloss.NewStyle().Foreground(Gray300).Render(truncate(reason, 25)),
						lipgloss.NewStyle().Foreground(Gray500).Render(debt.Date.Format("Jan 02")),
					)
					cardContent.WriteString(line)
				}
			}

			content.WriteString(card.Render(cardContent.String()))
			content.WriteString("\n")
			visibleIndex++
		}
	}

	// Summary bar at bottom
	stats := fmt.Sprintf("📥 Borrowed: %s  •  📤 Lent: %s  •  💰 Net: %s",
		MoneyNegativeStyle.Render(FormatMoneyPlain(data.TotalBorrowed(), m.config.Currency)),
		MoneyPositiveStyle.Render(FormatMoneyPlain(data.TotalLent(), m.config.Currency)),
		FormatAmount(data.TotalLent()-data.TotalBorrowed(), m.config.Currency),
	)

	// Help bar with styled keys
	helpBar := lipgloss.NewStyle().
		Background(Gray800).
		Foreground(Gray400).
		Padding(0, 2).
		Width(66).
		Render("a Add • s Settle • h History • i Details • Esc Back")

	// Combine everything
	result := lipgloss.JoinVertical(lipgloss.Left,
		header,
		content.String(),
		stats,
		helpBar,
	)

	return result
}

func (m *Model) updateDebtsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debts := m.storage.GetUnsettledDebts()

	// Build the same grouped structure as the view (by person name)
	type personGroup struct {
		name          string
		totalLent     float64
		totalBorrowed float64
		debts         []models.DebtTransaction
	}

	groupMap := make(map[string]*personGroup)
	var groupOrder []string

	for _, debt := range debts {
		key := debt.PersonName
		if _, exists := groupMap[key]; !exists {
			groupMap[key] = &personGroup{
				name:          debt.PersonName,
				totalLent:     0,
				totalBorrowed: 0,
				debts:         []models.DebtTransaction{},
			}
			groupOrder = append(groupOrder, key)
		}
		groupMap[key].debts = append(groupMap[key].debts, debt)
		if debt.Type == models.Lent {
			groupMap[key].totalLent += debt.Amount
		} else {
			groupMap[key].totalBorrowed += debt.Amount
		}
	}

	// Count only visible persons (net balance != 0)
	visibleCount := 0
	for _, key := range groupOrder {
		group := groupMap[key]
		netBalance := group.totalLent - group.totalBorrowed
		if netBalance != 0 {
			visibleCount++
		}
	}

	maxCursor := visibleCount - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	// Adjust cursor if it went out of bounds
	if m.cursor > maxCursor {
		m.cursor = maxCursor
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "a":
		m.currentView = ViewAddDebt
		m.initDebtInputs()
	case "s":
		// Open transaction selection view for settling
		if len(groupOrder) > 0 && m.cursor < len(groupOrder) {
			m.selectedPerson = groupOrder[m.cursor]
			m.currentView = ViewSelectTransactionToSettle
			m.cursor = 0
		}
	case "h":
		// View settlement history
		m.currentView = ViewDebtHistory
		m.cursor = 0
	case "i":
		// View person history (detailed view)
		if len(groupOrder) > 0 && m.cursor < len(groupOrder) {
			m.selectedPerson = groupOrder[m.cursor]
			m.currentView = ViewPersonHistory
			m.cursor = 0
		}
	}

	return m, nil
}

// Debt History view
func (m Model) viewDebtHistory() string {
	title := TitleStyle.Render("  Settlement History")

	settledDebts := m.storage.GetSettledDebts()

	var content string
	if len(settledDebts) == 0 {
		content = MutedStyle.Render("\n  No settled debts yet.\n")
	} else {
		// Sort by settled date (most recent first)
		for i := len(settledDebts) - 1; i >= 0; i-- {
			debt := settledDebts[i]

			// Skip if settled date is nil (shouldn't happen but be safe)
			if debt.SettledDate == nil {
				continue
			}

			cursor := "  "
			if len(settledDebts)-1-i == m.cursor {
				cursor = "▸ "
			}

			// Format the settlement info
			var typeStyle lipgloss.Style
			var typeLabel string
			if debt.Type == models.Lent {
				typeStyle = AmountPositiveStyle
				typeLabel = "Lent"
			} else {
				typeStyle = AmountNegativeStyle
				typeLabel = "Borrowed"
			}

			line := fmt.Sprintf("%s%s %s %s",
				cursor,
				debt.SettledDate.Format("2006-01-02"),
				SelectedMenuItemStyle.Render(debt.PersonName),
				typeStyle.Render(fmt.Sprintf("[%s]", typeLabel)),
			)
			content += line + "\n"

			// Show original amount and settlement details
			content += fmt.Sprintf("    Original: %s",
				FormatAmountPlain(debt.Amount, m.config.Currency),
			)

			// Show what was actually settled (might be different due to partial settlements)
			if debt.SettlementAmount > 0 && debt.SettlementAmount != debt.Amount {
				content += fmt.Sprintf(" | Settled: %s",
					FormatAmountPlain(debt.SettlementAmount, m.config.Currency),
				)
			}
			content += "\n"

			// Show description if any
			if debt.Description != "" {
				content += fmt.Sprintf("    Reason: %s\n",
					MutedStyle.Render(truncate(debt.Description, 40)),
				)
			}

			// Show settlement note
			if debt.SettlementNote != "" {
				content += fmt.Sprintf("    Note: %s\n",
					MutedStyle.Render(debt.SettlementNote),
				)
			} else {
				content += fmt.Sprintf("    Note: %s\n",
					MutedStyle.Render("(no note)"),
				)
			}

			content += "\n"
		}
	}

	help := HelpStyle.Render("\n  ↑/↓: Navigate • Esc: Back to debts")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateDebtHistoryView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	settledDebts := m.storage.GetSettledDebts()
	maxCursor := len(settledDebts) - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "esc", "q":
		m.currentView = ViewDebts
		m.cursor = 0
	}

	return m, nil
}

// Person History view - detailed view for a specific person
func (m Model) viewPersonHistory() string {
	title := TitleStyle.Render(fmt.Sprintf("  History: %s", m.selectedPerson))

	debts := m.storage.GetAllDebtsForPerson(m.selectedPerson)

	var content string
	if len(debts) == 0 {
		content = MutedStyle.Render(fmt.Sprintf("\n  No transactions with %s.\n", m.selectedPerson))
	} else {
		// Calculate totals
		var totalLent, totalBorrowed, totalSettledLent, totalSettledBorrowed float64
		var unsettledLent, unsettledBorrowed []models.DebtTransaction
		var settledTransactions []models.DebtTransaction

		for _, debt := range debts {
			if debt.Type == models.Lent {
				if debt.IsSettled {
					// Only count if settlement amount > 0
					if debt.SettlementAmount > 0 {
						totalSettledLent += debt.SettlementAmount
						settledTransactions = append(settledTransactions, debt)
					}
				} else {
					totalLent += debt.Amount
					unsettledLent = append(unsettledLent, debt)
				}
			} else {
				if debt.IsSettled {
					// Only count if settlement amount > 0
					if debt.SettlementAmount > 0 {
						totalSettledBorrowed += debt.SettlementAmount
						settledTransactions = append(settledTransactions, debt)
					}
				} else {
					totalBorrowed += debt.Amount
					unsettledBorrowed = append(unsettledBorrowed, debt)
				}
			}
		}

		// Summary section
		content = "\n"
		content += SelectedMenuItemStyle.Render("  Summary") + "\n"
		content += fmt.Sprintf("    Total Lent:       %s\n", AmountPositiveStyle.Render(FormatAmountPlain(totalLent, m.config.Currency)))
		content += fmt.Sprintf("    Total Borrowed:   %s\n", AmountNegativeStyle.Render(FormatAmountPlain(totalBorrowed, m.config.Currency)))
		content += fmt.Sprintf("    Settled (Lent):   %s\n", FormatAmountPlain(totalSettledLent, m.config.Currency))
		content += fmt.Sprintf("    Settled (Borrow): %s\n", FormatAmountPlain(totalSettledBorrowed, m.config.Currency))
		netBalance := (totalLent + totalSettledLent) - (totalBorrowed + totalSettledBorrowed)
		content += fmt.Sprintf("    Net Balance:      %s\n\n", FormatAmount(netBalance, m.config.Currency))

		// Unsettled Transactions
		if len(unsettledLent) > 0 || len(unsettledBorrowed) > 0 {
			content += SelectedMenuItemStyle.Render("  Active Transactions") + "\n"

			if len(unsettledLent) > 0 {
				content += AmountPositiveStyle.Render("    Lent:") + "\n"
				for _, debt := range unsettledLent {
					content += fmt.Sprintf("      • %s  %s  %s\n",
						debt.Date.Format("2006-01-02"),
						FormatAmountPlain(debt.Amount, m.config.Currency),
						MutedStyle.Render(truncate(debt.Description, 30)),
					)
				}
			}

			if len(unsettledBorrowed) > 0 {
				content += AmountNegativeStyle.Render("    Borrowed:") + "\n"
				for _, debt := range unsettledBorrowed {
					content += fmt.Sprintf("      • %s  %s  %s\n",
						debt.Date.Format("2006-01-02"),
						FormatAmountPlain(debt.Amount, m.config.Currency),
						MutedStyle.Render(truncate(debt.Description, 30)),
					)
				}
			}
			content += "\n"
		}

		// Settled Transactions
		if len(settledTransactions) > 0 {
			content += SelectedMenuItemStyle.Render("  Settlement History") + "\n"

			// Sort by settled date (most recent first)
			for i := len(settledTransactions) - 1; i >= 0; i-- {
				debt := settledTransactions[i]
				if debt.SettledDate == nil {
					continue
				}

				var typeLabel string
				if debt.Type == models.Lent {
					typeLabel = "Lent"
				} else {
					typeLabel = "Borrowed"
				}

				content += fmt.Sprintf("    ✓ %s  %s  %s\n",
					debt.SettledDate.Format("2006-01-02"),
					FormatAmountPlain(debt.SettlementAmount, m.config.Currency),
					typeLabel,
				)

				if debt.Description != "" {
					content += fmt.Sprintf("      Original: %s\n",
						MutedStyle.Render(debt.Description),
					)
				}

				if debt.SettlementNote != "" {
					content += fmt.Sprintf("      Note: %s\n",
						MutedStyle.Render(debt.SettlementNote),
					)
				}
				content += "\n"
			}
		}
	}

	help := HelpStyle.Render("\n  Esc: Back to debts")

	return BoxStyle.Render(title + content + help)
}

func (m *Model) updatePersonHistoryView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.currentView = ViewDebts
		m.cursor = 0
		// Don't clear selectedPerson so we return to the same person in the list
	}

	return m, nil
}

// Select Transaction to Settle view
func (m Model) viewSelectTransactionToSettle() string {
	title := TitleStyle.Render("  Select Transaction to Settle")

	debts := m.storage.GetUnsettledDebtsForPerson(m.selectedPerson)

	var content string
	if len(debts) == 0 {
		content = MutedStyle.Render("\n  No unsettled debts for " + m.selectedPerson + ".\n")
	} else {
		content = "\n  " + SelectedMenuItemStyle.Render(m.selectedPerson) + "\n\n"

		for i, debt := range debts {
			cursor := "  "
			if i == m.cursor {
				cursor = "▸ "
			}

			var typeStyle lipgloss.Style
			var typeLabel string
			if debt.Type == models.Lent {
				typeStyle = AmountPositiveStyle
				typeLabel = "Lent"
			} else {
				typeStyle = AmountNegativeStyle
				typeLabel = "Borrowed"
			}

			line := fmt.Sprintf("%s%s  %s  %s",
				cursor,
				debt.Date.Format("2006-01-02"),
				typeStyle.Render(fmt.Sprintf("[%s]", typeLabel)),
				FormatAmountPlain(debt.Amount, m.config.Currency),
			)
			content += line + "\n"

			// Show description
			desc := debt.Description
			if desc == "" {
				desc = "(no description)"
			}
			content += fmt.Sprintf("    %s\n", MutedStyle.Render(truncate(desc, 40)))
			content += "\n"
		}
	}

	help := HelpStyle.Render("\n  ↑/↓: Navigate • Enter: Select • Esc: Back")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateSelectTransactionToSettleView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	debts := m.storage.GetUnsettledDebtsForPerson(m.selectedPerson)
	maxCursor := len(debts) - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "enter":
		if len(debts) > 0 && m.cursor < len(debts) {
			m.selectedID = debts[m.cursor].ID
			m.currentView = ViewSettleDebt
			m.initSettleDebtInputs()
		}
	case "esc":
		m.currentView = ViewDebts
		m.cursor = 0
		m.selectedPerson = ""
	}

	return m, nil
}

func (m *Model) initDebtInputs() {
	m.inputs = make([]textinput.Model, 5)

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Type (borrowed/lent)"
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Person Name"

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "Amount"

	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "Description"

	m.inputs[4] = textinput.New()
	m.inputs[4].Placeholder = "Transaction Date (YYYY-MM-DD)"
	m.inputs[4].SetValue(time.Now().Format("2006-01-02"))

	m.focusIndex = 0
}

func (m Model) viewAddDebt() string {
	title := TitleStyle.Render("  Add Debt Transaction")

	var content string
	labels := []string{"Type:", "Person:", "Amount:", "Description:", "Date:"}
	hints := []string{
		"Options: borrowed, lent",
		"",
		"",
		"",
		"Date when borrowed/lent (YYYY-MM-DD)",
	}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		} else {
			content += MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		}
	}

	help := HelpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateAddDebtView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		txType := models.TransactionType(m.inputs[0].Value())
		if txType != models.Borrowed && txType != models.Lent {
			m.message = "Type must be 'borrowed' or 'lent'"
			m.messageType = "error"
			return m, nil
		}

		personName := m.inputs[1].Value()
		if personName == "" {
			m.message = "Person name is required"
			m.messageType = "error"
			return m, nil
		}

		amount, err := strconv.ParseFloat(m.inputs[2].Value(), 64)
		if err != nil {
			m.message = "Invalid amount"
			m.messageType = "error"
			return m, nil
		}

		description := m.inputs[3].Value()

		// Parse transaction date
		dateStr := m.inputs[4].Value()
		if dateStr == "" {
			m.message = "Transaction date is required"
			m.messageType = "error"
			return m, nil
		}
		transactionDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			m.message = "Invalid date format. Use YYYY-MM-DD"
			m.messageType = "error"
			return m, nil
		}

		_, err = m.storage.AddDebtTransaction(txType, personName, amount, description, transactionDate, nil)
		if err != nil {
			m.message = "Error saving: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Debt transaction added!"
		m.messageType = "success"
		m.currentView = ViewDebts
		m.inputs = nil
		m.cursor = 0
		return m, nil
	case "+":
		if m.focusIndex == 2 && len(m.inputs) > 0 {
			currentValue := m.inputs[2].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[2].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

// Settle Debt functions
func (m *Model) initSettleDebtInputs() {
	// Get the selected transaction
	var selectedTx *models.DebtTransaction
	for _, tx := range m.storage.GetDebtTransactions() {
		if tx.ID == m.selectedID {
			selectedTx = &tx
			break
		}
	}

	var maxAmount float64
	if selectedTx != nil {
		maxAmount = selectedTx.Amount
	}

	m.inputs = make([]textinput.Model, 2)
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = fmt.Sprintf("Amount to settle (max: %.2f)", maxAmount)
	m.inputs[0].SetValue(fmt.Sprintf("%.2f", maxAmount))
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Settlement note (e.g., Cash payment, Bank transfer, etc.)"

	m.focusIndex = 0
}

func (m Model) viewSettleDebt() string {
	title := TitleStyle.Render("  Settle Transaction")

	// Get the selected transaction details
	var selectedTx *models.DebtTransaction
	for _, tx := range m.storage.GetDebtTransactions() {
		if tx.ID == m.selectedID {
			selectedTx = &tx
			break
		}
	}

	var content string
	if selectedTx != nil {
		var typeStyle lipgloss.Style
		var typeLabel string
		if selectedTx.Type == models.Lent {
			typeStyle = AmountPositiveStyle
			typeLabel = "Lent"
		} else {
			typeStyle = AmountNegativeStyle
			typeLabel = "Borrowed"
		}

		content += fmt.Sprintf("  %s: %s\n", SelectedMenuItemStyle.Render("Person"), selectedTx.PersonName)
		content += fmt.Sprintf("  %s: %s\n", SelectedMenuItemStyle.Render("Type"), typeStyle.Render(typeLabel))
		content += fmt.Sprintf("  %s: %s\n", SelectedMenuItemStyle.Render("Original Amount"), FormatAmountPlain(selectedTx.Amount, m.config.Currency))
		content += fmt.Sprintf("  %s: %s\n", SelectedMenuItemStyle.Render("Date"), selectedTx.Date.Format("2006-01-02"))
		if selectedTx.Description != "" {
			content += fmt.Sprintf("  %s: %s\n", SelectedMenuItemStyle.Render("Description"), selectedTx.Description)
		}
		content += "\n"
	}

	if len(m.inputs) > 0 {
		// Amount field
		if m.focusIndex == 0 {
			content += "  " + SelectedMenuItemStyle.Render("▸ Amount to settle:") + "\n"
			content += "  " + FocusedInputStyle.Render(m.inputs[0].View()) + "\n"
		} else {
			content += "  " + MenuItemStyle.Render("  Amount to settle:") + "\n"
			content += "  " + InputStyle.Render(m.inputs[0].View()) + "\n"
		}
		content += "  " + MutedStyle.Render("Enter amount (defaults to full amount)") + "\n\n"

		// Settlement note field
		if m.focusIndex == 1 {
			content += "  " + SelectedMenuItemStyle.Render("▸ Settlement note:") + "\n"
			content += "  " + FocusedInputStyle.Render(m.inputs[1].View()) + "\n"
		} else {
			content += "  " + MenuItemStyle.Render("  Settlement note:") + "\n"
			content += "  " + InputStyle.Render(m.inputs[1].View()) + "\n"
		}
		content += "  " + MutedStyle.Render("Why/How it was settled (e.g., Cash, Bank transfer, UPI)") + "\n"
	}

	help := HelpStyle.Render("\n  Tab: Next field • Enter: Confirm • Esc: Cancel")

	return BoxStyle.Render(title + "\n\n" + content + help)
}

func (m *Model) updateSettleDebtView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
		return m, nil
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
		return m, nil
	case "enter":
		// Get the selected transaction to find max amount
		var selectedTx *models.DebtTransaction
		for _, tx := range m.storage.GetDebtTransactions() {
			if tx.ID == m.selectedID {
				selectedTx = &tx
				break
			}
		}

		if selectedTx == nil {
			m.message = "Transaction not found"
			m.messageType = "error"
			return m, nil
		}

		maxAmount := selectedTx.Amount
		var amountToSettle float64

		if m.inputs[0].Value() == "" {
			amountToSettle = maxAmount
		} else {
			var err error
			amountToSettle, err = strconv.ParseFloat(m.inputs[0].Value(), 64)
			if err != nil {
				m.message = "Invalid amount"
				m.messageType = "error"
				return m, nil
			}
			if amountToSettle <= 0 {
				m.message = "Amount must be positive"
				m.messageType = "error"
				return m, nil
			}
			if amountToSettle > maxAmount {
				m.message = fmt.Sprintf("Amount exceeds transaction amount (max: %.2f)", maxAmount)
				m.messageType = "error"
				return m, nil
			}
		}

		settlementNote := m.inputs[1].Value()
		err := m.storage.SettleTransaction(m.selectedID, amountToSettle, settlementNote)
		if err != nil {
			m.message = "Error settling: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = fmt.Sprintf("Settled %s %.2f!", m.config.Currency, amountToSettle)
		m.messageType = "success"
		m.currentView = ViewDebts
		m.inputs = nil
		m.selectedID = ""
		m.selectedPerson = ""
		m.cursor = 0
		return m, nil
	case "+":
		if m.focusIndex == 0 && len(m.inputs) > 0 {
			currentValue := m.inputs[0].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[0].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

// Net Worth view
func (m Model) viewNetWorth() string {
	title := TitleStyle.Render("  My Net Worth")

	investments := m.storage.GetInvestments()
	data := m.storage.GetData()

	var content string
	if len(investments) == 0 {
		content = MutedStyle.Render("\n  No investments recorded yet.\n")
	} else {
		content = "\n"
		for i, inv := range investments {
			cursor := "  "
			if i == m.cursor {
				cursor = "▸ "
			}
			gain := inv.CurrentValue - inv.InvestedAmount
			gainPct := float64(0)
			if inv.InvestedAmount > 0 {
				gainPct = (gain / inv.InvestedAmount) * 100
			}
			line := fmt.Sprintf("%s[%s] %s  %s  %s (%.1f%%)",
				cursor,
				TableCellStyle.Width(12).Render(string(inv.Type)),
				TableCellStyle.Width(20).Render(truncate(inv.Name, 20)),
				FormatAmountPlain(inv.CurrentValue, m.config.Currency),
				FormatAmount(gain, ""),
				gainPct,
			)
			content += line + "\n"
		}
	}

	// Summary
	netWorth := data.NetWorth()
	stats := fmt.Sprintf("\n  Total Net Worth: %s", FormatAmountPlain(netWorth, m.config.Currency))

	help := HelpStyle.Render("\n  a: Add investment • u: Update value • d: Delete • Esc: Back")

	return BoxStyle.Render(title + content + stats + help)
}

func (m *Model) updateNetWorthView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	investments := m.storage.GetInvestments()
	maxCursor := len(investments) - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "a":
		m.currentView = ViewAddInvestment
		m.initInvestmentInputs()
	case "d":
		if len(investments) > 0 && m.cursor < len(investments) {
			m.currentView = ViewConfirmDelete
			m.inputs = nil
		}
	case "u":
		if len(investments) > 0 && m.cursor < len(investments) {
			m.selectedID = investments[m.cursor].ID
			m.currentView = ViewUpdateInvestment
			m.inputs = make([]textinput.Model, 2)
			m.inputs[0] = textinput.New()
			m.inputs[0].Placeholder = "New invested amount"
			m.inputs[0].SetValue(fmt.Sprintf("%.2f", investments[m.cursor].InvestedAmount))
			m.inputs[1] = textinput.New()
			m.inputs[1].Placeholder = "New current value"
			m.inputs[1].SetValue(fmt.Sprintf("%.2f", investments[m.cursor].CurrentValue))
			m.inputs[0].Focus()
			m.focusIndex = 0
		}
	}

	return m, nil
}

func (m *Model) initInvestmentInputs() {
	m.inputs = make([]textinput.Model, 6)

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Type (stocks/mutual_funds/gold/silver/fixed_deposit/ppf/crypto/other)"
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Name (e.g., HDFC Bank, SBI Bluechip)"

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "Invested Amount"

	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "Current Value"

	m.inputs[4] = textinput.New()
	m.inputs[4].Placeholder = "Units (optional)"

	m.inputs[5] = textinput.New()
	m.inputs[5].Placeholder = "Purchase Date (YYYY-MM-DD)"

	m.focusIndex = 0
}

func (m Model) viewAddInvestment() string {
	title := TitleStyle.Render("  Add Investment")

	var content string
	labels := []string{"Type:", "Name:", "Invested:", "Current Value:", "Units:", "Purchase Date:"}
	hints := []string{
		"Options: stocks, mutual_funds, gold, silver, fixed_deposit, ppf, crypto, real_estate, other",
		"e.g., HDFC Bank, SBI Bluechip, Gold 24K",
		"",
		"",
		"(optional)",
		"Format: YYYY-MM-DD",
	}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		} else {
			content += MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		}
	}

	help := HelpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m Model) viewUpdateInvestment() string {
	title := TitleStyle.Render("  Update Investment Value")

	var content string
	content += "\n"

	labels := []string{"New invested amount:", "New current value:"}
	hints := []string{"Enter the new invested amount", "Enter the new current value"}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += "  " + SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		} else {
			content += "  " + MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n"
			if hints[i] != "" {
				content += "  " + MutedStyle.Render(hints[i]) + "\n"
			}
			content += "\n"
		}
	}

	help := HelpStyle.Render("\n  Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m Model) viewConfirmDelete() string {
	title := TitleStyle.Render("  Confirm Delete")

	var content string
	content += "\n  Are you sure you want to delete this investment?\n\n"
	content += "  This action cannot be undone.\n"

	help := HelpStyle.Render("\n  Enter: Yes, delete • Esc: Cancel")

	return BoxStyle.Render(title + content + help)
}

func (m *Model) updateAddInvestmentView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		invType := models.InvestmentType(m.inputs[0].Value())
		name := m.inputs[1].Value()
		if name == "" {
			m.message = "Name is required"
			m.messageType = "error"
			return m, nil
		}

		invested, err := strconv.ParseFloat(m.inputs[2].Value(), 64)
		if err != nil {
			m.message = "Invalid invested amount"
			m.messageType = "error"
			return m, nil
		}

		current, err := strconv.ParseFloat(m.inputs[3].Value(), 64)
		if err != nil {
			m.message = "Invalid current value"
			m.messageType = "error"
			return m, nil
		}

		var units float64
		if m.inputs[4].Value() != "" {
			units, _ = strconv.ParseFloat(m.inputs[4].Value(), 64)
		}

		purchaseDate := time.Now()
		if m.inputs[5].Value() != "" {
			purchaseDate, err = time.Parse("2006-01-02", m.inputs[5].Value())
			if err != nil {
				m.message = "Invalid date format"
				m.messageType = "error"
				return m, nil
			}
		}

		_, err = m.storage.AddInvestment(invType, name, invested, current, units, purchaseDate, "")
		if err != nil {
			m.message = "Error saving: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Investment added!"
		m.messageType = "success"
		m.currentView = ViewNetWorth
		m.inputs = nil
		m.cursor = 0
		return m, nil
	case "+":
		if (m.focusIndex == 2 || m.focusIndex == 3) && len(m.inputs) > 0 {
			currentValue := m.inputs[m.focusIndex].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[m.focusIndex].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) updateUpdateInvestmentView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		if m.inputs[0].Value() == "" || m.inputs[1].Value() == "" {
			m.message = "Both values are required"
			m.messageType = "error"
			return m, nil
		}

		investedAmount, err := strconv.ParseFloat(m.inputs[0].Value(), 64)
		if err != nil {
			m.message = "Invalid invested amount"
			m.messageType = "error"
			return m, nil
		}

		currentValue, err := strconv.ParseFloat(m.inputs[1].Value(), 64)
		if err != nil {
			m.message = "Invalid current value"
			m.messageType = "error"
			return m, nil
		}

		if investedAmount < 0 || currentValue < 0 {
			m.message = "Values must be positive"
			m.messageType = "error"
			return m, nil
		}

		err = m.storage.UpdateInvestment(m.selectedID, investedAmount, currentValue)
		if err != nil {
			m.message = "Error updating: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Investment updated!"
		m.messageType = "success"
		m.currentView = ViewNetWorth
		m.inputs = nil
		m.selectedID = ""
		m.cursor = 0
		return m, nil
	case "+":
		if (m.focusIndex == 0 || m.focusIndex == 1) && len(m.inputs) > 0 {
			currentValue := m.inputs[m.focusIndex].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[m.focusIndex].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) updateConfirmDeleteView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		investments := m.storage.GetInvestments()
		if len(investments) > 0 && m.cursor < len(investments) {
			m.storage.DeleteInvestment(investments[m.cursor].ID)
			m.message = "Investment deleted"
			m.messageType = "success"
		}
		m.currentView = ViewNetWorth
		m.inputs = nil
		m.cursor = 0
		return m, nil
	}

	return m, nil
}

// Savings view
func (m Model) viewSavings() string {
	title := TitleStyle.Render("  Savings Goals")

	targets := m.storage.GetSavingsTargets()

	var content string
	if len(targets) == 0 {
		content = MutedStyle.Render("\n  No savings goals set yet.\n")
	} else {
		content = "\n"
		for i, target := range targets {
			cursor := "  "
			if i == m.cursor {
				cursor = "▸ "
			}
			status := "Active"
			if target.IsCompleted {
				status = "Done!"
			}
			line := fmt.Sprintf("%s%s\n    %s / %s  [%s]\n    %s  Due: %s\n",
				cursor,
				SelectedMenuItemStyle.Render(target.ProductName),
				FormatAmountPlain(target.CurrentAmount, m.config.Currency),
				FormatAmountPlain(target.TargetAmount, m.config.Currency),
				status,
				ProgressBar(target.CurrentAmount, target.TargetAmount, 20),
				target.TargetDate.Format("2006-01-02"),
			)
			content += line
		}
	}

	help := HelpStyle.Render("\n  a: Add goal • c: Add contribution • d: Delete • Esc: Back")

	return BoxStyle.Render(title + content + help)
}

func (m *Model) updateSavingsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	targets := m.storage.GetSavingsTargets()
	maxCursor := len(targets) - 1
	if maxCursor < 0 {
		maxCursor = 0
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < maxCursor {
			m.cursor++
		}
	case "a":
		m.currentView = ViewAddSavingsTarget
		m.initSavingsTargetInputs()
	case "c":
		if len(targets) > 0 && m.cursor < len(targets) {
			m.selectedID = targets[m.cursor].ID
			m.currentView = ViewAddContribution
			m.initContributionInputs()
		}
	case "d":
		if len(targets) > 0 && m.cursor < len(targets) {
			m.storage.DeleteSavingsTarget(targets[m.cursor].ID)
			m.message = "Goal deleted"
			m.messageType = "success"
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}

	return m, nil
}

func (m *Model) initSavingsTargetInputs() {
	m.inputs = make([]textinput.Model, 4)

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Product Name (e.g., iPhone 16, MacBook Pro)"
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Target Amount"

	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "Target Date (YYYY-MM-DD)"

	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "Description (optional)"

	m.focusIndex = 0
}

func (m Model) viewAddSavingsTarget() string {
	title := TitleStyle.Render("  Add Savings Goal")

	var content string
	labels := []string{"Product:", "Target Amount:", "Target Date:", "Description:"}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n\n"
		} else {
			content += MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n\n"
		}
	}

	help := HelpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateAddSavingsTargetView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		productName := m.inputs[0].Value()
		if productName == "" {
			m.message = "Product name is required"
			m.messageType = "error"
			return m, nil
		}

		targetAmount, err := strconv.ParseFloat(m.inputs[1].Value(), 64)
		if err != nil {
			m.message = "Invalid target amount"
			m.messageType = "error"
			return m, nil
		}

		targetDate, err := time.Parse("2006-01-02", m.inputs[2].Value())
		if err != nil {
			m.message = "Invalid date format (use YYYY-MM-DD)"
			m.messageType = "error"
			return m, nil
		}

		description := m.inputs[3].Value()

		_, err = m.storage.AddSavingsTarget(productName, targetAmount, targetDate, description)
		if err != nil {
			m.message = "Error saving: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Savings goal created!"
		m.messageType = "success"
		m.currentView = ViewSavings
		m.inputs = nil
		m.cursor = 0
		return m, nil
	case "+":
		if m.focusIndex == 1 && len(m.inputs) > 0 {
			currentValue := m.inputs[1].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[1].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) initContributionInputs() {
	m.inputs = make([]textinput.Model, 2)

	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Amount"
	m.inputs[0].Focus()

	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "Notes (optional)"

	m.focusIndex = 0
}

func (m Model) viewAddContribution() string {
	title := TitleStyle.Render("  Add Contribution")

	var content string
	labels := []string{"Amount:", "Notes:"}

	for i, input := range m.inputs {
		label := labels[i]
		if i == m.focusIndex {
			content += SelectedMenuItemStyle.Render("▸ "+label) + "\n"
			content += "  " + FocusedInputStyle.Render(input.View()) + "\n\n"
		} else {
			content += MenuItemStyle.Render("  "+label) + "\n"
			content += "  " + InputStyle.Render(input.View()) + "\n\n"
		}
	}

	help := HelpStyle.Render("Tab: Next field • Enter: Save • Esc: Cancel")

	return BoxStyle.Render(title + "\n" + content + help)
}

func (m *Model) updateAddContributionView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "down":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
		m.inputs[m.focusIndex].Focus()
	case "shift+tab", "up":
		m.inputs[m.focusIndex].Blur()
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs) - 1
		}
		m.inputs[m.focusIndex].Focus()
	case "enter":
		amount, err := strconv.ParseFloat(m.inputs[0].Value(), 64)
		if err != nil {
			m.message = "Invalid amount"
			m.messageType = "error"
			return m, nil
		}

		notes := m.inputs[1].Value()

		_, err = m.storage.AddSavingsContribution(m.selectedID, amount, notes)
		if err != nil {
			m.message = "Error saving: " + err.Error()
			m.messageType = "error"
			return m, nil
		}

		m.message = "Contribution added!"
		m.messageType = "success"
		m.currentView = ViewSavings
		m.inputs = nil
		m.selectedID = ""
		m.cursor = 0
		return m, nil
	case "+":
		if m.focusIndex == 0 && len(m.inputs) > 0 {
			currentValue := m.inputs[0].Value()
			calculatedValue, success := tryCalculateAmount(currentValue)
			if success {
				m.inputs[0].SetValue(calculatedValue)
				m.message = "Calculated: " + calculatedValue
				m.messageType = "info"
			}
		}
	}

	if len(m.inputs) > 0 && m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}
	return m, nil
}

// Stats view
func (m Model) viewStats() string {
	title := TitleStyle.Render("  Stats & Dashboard")

	data := m.storage.GetData()
	now := time.Now()

	// Net Worth
	netWorth := data.NetWorth()

	// Debts
	totalBorrowed := data.TotalBorrowed()
	totalLent := data.TotalLent()

	// Expenses
	monthlyExpenses := data.MonthlyExpenses(now.Year(), now.Month())
	var totalExpenses float64
	for _, e := range data.Expenses {
		totalExpenses += e.Amount
	}

	// Savings
	var activeSavings, completedSavings int
	var totalSavingsTarget, totalSaved float64
	for _, t := range data.SavingsTargets {
		if t.IsCompleted {
			completedSavings++
		} else {
			activeSavings++
		}
		totalSavingsTarget += t.TargetAmount
		totalSaved += t.CurrentAmount
	}

	content := fmt.Sprintf(`
  %s
  ──────────────────────────
  Total Net Worth:     %s

  %s
  ──────────────────────────
  Total Borrowed:      %s
  Total Lent:          %s
  Net Position:        %s

  %s
  ──────────────────────────
  This Month:          %s
  All Time:            %s

  %s
  ──────────────────────────
  Active Goals:        %d
  Completed Goals:     %d
  Total Target:        %s
  Total Saved:         %s
  Progress:            %s
`,
		SelectedMenuItemStyle.Render("NET WORTH"),
		FormatAmountPlain(netWorth, m.config.Currency),
		SelectedMenuItemStyle.Render("DEBTS"),
		FormatAmountPlain(totalBorrowed, m.config.Currency),
		FormatAmountPlain(totalLent, m.config.Currency),
		FormatAmount(totalLent-totalBorrowed, m.config.Currency),
		SelectedMenuItemStyle.Render("EXPENSES"),
		FormatAmountPlain(monthlyExpenses, m.config.Currency),
		FormatAmountPlain(totalExpenses, m.config.Currency),
		SelectedMenuItemStyle.Render("SAVINGS GOALS"),
		activeSavings,
		completedSavings,
		FormatAmountPlain(totalSavingsTarget, m.config.Currency),
		FormatAmountPlain(totalSaved, m.config.Currency),
		ProgressBar(totalSaved, totalSavingsTarget, 20),
	)

	help := HelpStyle.Render("\n  Esc: Back to main menu")

	return BoxStyle.Render(title + content + help)
}

func (m *Model) updateStatsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Stats view is read-only, just handle navigation
	return m, nil
}

// Helper functions
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func (m *Model) setMessage(msg, msgType string) {
	m.message = msg
	m.messageType = msgType
}

func evaluateMathExpression(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, fmt.Errorf("empty expression")
	}

	parts := strings.Split(expr, "+")
	if len(parts) > 1 {
		var result float64
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			result += val
		}
		return result, nil
	}

	parts = strings.Split(expr, "-")
	if len(parts) > 1 {
		val0, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, err
		}
		result := val0
		for _, part := range parts[1:] {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			result -= val
		}
		return result, nil
	}

	parts = strings.Split(expr, "*")
	if len(parts) > 1 {
		result := 1.0
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			result *= val
		}
		return result, nil
	}

	parts = strings.Split(expr, "/")
	if len(parts) > 1 {
		val0, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, err
		}
		result := val0
		for _, part := range parts[1:] {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			if val == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			result /= val
		}
		return result, nil
	}

	return strconv.ParseFloat(expr, 64)
}

func tryCalculateAmount(value string) (string, bool) {
	cleanValue := strings.TrimSpace(value)

	if strings.ContainsAny(cleanValue, "+-*/") {
		for _, char := range cleanValue {
			if !unicode.IsDigit(char) && char != '.' && char != ' ' && char != '+' && char != '-' && char != '*' && char != '/' {
				return value, false
			}
		}

		result, err := evaluateMathExpression(cleanValue)
		if err == nil {
			return strconv.FormatFloat(result, 'f', -1, 64), true
		}
	}

	return value, false
}
