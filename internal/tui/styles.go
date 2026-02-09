package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	Primary       = lipgloss.Color("#7C3AED")
	Secondary     = lipgloss.Color("#10B981")
	Accent        = lipgloss.Color("#F59E0B")
	Danger        = lipgloss.Color("#EF4444")
	Muted         = lipgloss.Color("#6B7280")
	Background    = lipgloss.Color("#1F2937")
	Surface       = lipgloss.Color("#374151")
	TextPrimary   = lipgloss.Color("#F9FAFB")
	TextSecondary = lipgloss.Color("#9CA3AF")
)

// Styles
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(TextSecondary).
			MarginBottom(1)

	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(Primary).
				Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Accent)

	MutedStyle = lipgloss.NewStyle().
			Foreground(Muted)

	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder())

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)

	AmountPositiveStyle = lipgloss.NewStyle().
				Foreground(Secondary)

	AmountNegativeStyle = lipgloss.NewStyle().
				Foreground(Danger)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(Primary)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			MarginTop(1)

	InputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	FocusedInputStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(Secondary).
				Padding(0, 1)
)

// FormatAmount formats amount with color based on positive/negative
func FormatAmount(amount float64, currency string) string {
	if amount >= 0 {
		return AmountPositiveStyle.Render(currency + " " + formatFloat(amount))
	}
	return AmountNegativeStyle.Render(currency + " " + formatFloat(amount))
}

// FormatAmountPlain formats amount without styling
func FormatAmountPlain(amount float64, currency string) string {
	return currency + " " + formatFloat(amount)
}

func formatFloat(f float64) string {
	return lipgloss.NewStyle().Render(floatToString(f))
}

func floatToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

// ProgressBar creates a visual progress bar
func ProgressBar(current, total float64, width int) string {
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

	return ProgressBarStyle.Render(bar) + MutedStyle.Render(fmt.Sprintf(" %.1f%%", pct*100))
}
