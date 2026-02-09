package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Modern Color Palette - Tailwind-inspired
var (
	// Primary Colors
	Primary    = lipgloss.Color("#8B5CF6") // Violet 500
	Primary600 = lipgloss.Color("#7C3AED") // Violet 600
	Primary400 = lipgloss.Color("#A78BFA") // Violet 400

	// Secondary Colors
	Secondary    = lipgloss.Color("#10B981") // Emerald 500
	Secondary600 = lipgloss.Color("#059669") // Emerald 600

	// Accent Colors
	Accent    = lipgloss.Color("#F59E0B") // Amber 500
	Accent600 = lipgloss.Color("#D97706") // Amber 600

	// Status Colors
	Success = lipgloss.Color("#22C55E") // Green 500
	Warning = lipgloss.Color("#EAB308") // Yellow 500
	Danger  = lipgloss.Color("#EF4444") // Red 500
	Info    = lipgloss.Color("#3B82F6") // Blue 500

	// Neutral Colors
	Gray50  = lipgloss.Color("#F9FAFB")
	Gray100 = lipgloss.Color("#F3F4F6")
	Gray200 = lipgloss.Color("#E5E7EB")
	Gray300 = lipgloss.Color("#D1D5DB")
	Gray400 = lipgloss.Color("#9CA3AF")
	Gray500 = lipgloss.Color("#6B7280")
	Gray600 = lipgloss.Color("#4B5563")
	Gray700 = lipgloss.Color("#374151")
	Gray800 = lipgloss.Color("#1F2937")
	Gray900 = lipgloss.Color("#111827")
	Gray950 = lipgloss.Color("#030712")
)

// Enhanced Core Styles
var (
	// Base container with shadow effect
	BaseStyle = lipgloss.NewStyle().
			Background(Gray950)

	// App Header
	AppTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Background(lipgloss.Color("#1E1B4B")).
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			BorderBottom(true).
			Align(lipgloss.Center)

	AppSubtitleStyle = lipgloss.NewStyle().
				Foreground(Gray400).
				Italic(true).
				Align(lipgloss.Center)

	// Cards
	CardStyle = lipgloss.NewStyle().
			Background(Gray900).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gray700).
			Padding(1, 2)

	CardSelectedStyle = lipgloss.NewStyle().
				Background(Gray800).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(1, 2)

	CardHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Gray100).
			Background(Gray800).
			Padding(0, 1).
			Width(30)

	// Menu Items
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			PaddingLeft(2).
			PaddingTop(0).
			PaddingBottom(0)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Background(lipgloss.Color("#2E1065")).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1).
				BorderStyle(lipgloss.ThickBorder()).
				BorderLeft(true).
				BorderForeground(Primary)

	// Section Headers
	SectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary400).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(Gray700).
				PaddingBottom(1).
				MarginTop(1)

	// Money/Balance Styles
	MoneyPositiveStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true)

	MoneyNegativeStyle = lipgloss.NewStyle().
				Foreground(Danger).
				Bold(true)

	MoneyNeutralStyle = lipgloss.NewStyle().
				Foreground(Gray300)

	// Status badges
	BadgeSuccessStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#064E3B")).
				Foreground(Success).
				Padding(0, 1).
				Bold(true)

	BadgeWarningStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#713F12")).
				Foreground(Accent).
				Padding(0, 1).
				Bold(true)

	BadgeDangerStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7F1D1D")).
				Foreground(Danger).
				Padding(0, 1).
				Bold(true)

	BadgeInfoStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1E3A8A")).
			Foreground(Info).
			Padding(0, 1).
			Bold(true)

	// Table Styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Primary).
				Background(Gray800).
				Padding(0, 1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true).
				BorderForeground(Primary)

	TableCellStyle = lipgloss.NewStyle().
			Foreground(Gray300).
			Padding(0, 1)

	TableCellAltStyle = lipgloss.NewStyle().
				Foreground(Gray300).
				Background(Gray900).
				Padding(0, 1)

	// Form/Input Styles
	InputStyle = lipgloss.NewStyle().
			Foreground(Gray200).
			Background(Gray800).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Gray600).
			Padding(0, 1).
			Width(40)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(Gray100).
				Background(Gray800).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(0, 1).
				Width(40)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Gray400).
			Bold(true).
			MarginBottom(0)

	// Help/Footer
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Background(Gray800).
			Padding(0, 1)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Gray500)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Background(Gray900).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(Gray700).
			Padding(0, 1).
			Width(80)

	// Message Styles
	SuccessMessageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#064E3B")).
				Foreground(Success).
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Success)

	ErrorMessageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#7F1D1D")).
				Foreground(Danger).
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Danger)

	InfoMessageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1E3A8A")).
				Foreground(Info).
				Padding(1, 2).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Info)

	// Stats/Numbers
	StatNumberStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Gray100).
			Padding(0, 1)

	StatLabelStyle = lipgloss.NewStyle().
			Foreground(Gray500).
			Padding(0, 1)

	// Old aliases for compatibility (keeping old code working)
	TitleStyle          = AppTitleStyle
	SubtitleStyle       = AppSubtitleStyle
	BoxStyle            = CardStyle
	SuccessStyle        = lipgloss.NewStyle().Foreground(Success)
	ErrorStyle          = lipgloss.NewStyle().Foreground(Danger)
	WarningStyle        = lipgloss.NewStyle().Foreground(Accent)
	MutedStyle          = lipgloss.NewStyle().Foreground(Gray500)
	AmountPositiveStyle = MoneyPositiveStyle
	AmountNegativeStyle = MoneyNegativeStyle
	ProgressBarStyle    = lipgloss.NewStyle().Foreground(Primary)
	HelpStyle           = FooterStyle
)

// Utility Functions

// FormatMoney formats an amount with appropriate styling
func FormatMoney(amount float64, currency string) string {
	formatted := fmt.Sprintf("%s %.2f", currency, amount)
	if amount > 0 {
		return MoneyPositiveStyle.Render(formatted)
	} else if amount < 0 {
		return MoneyNegativeStyle.Render(formatted)
	}
	return MoneyNeutralStyle.Render(formatted)
}

// FormatMoneyPlain formats money without styling
func FormatMoneyPlain(amount float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, amount)
}

// RenderBadge creates a status badge
func RenderBadge(text string, status string) string {
	switch status {
	case "success":
		return BadgeSuccessStyle.Render(text)
	case "warning":
		return BadgeWarningStyle.Render(text)
	case "danger":
		return BadgeDangerStyle.Render(text)
	case "info":
		return BadgeInfoStyle.Render(text)
	default:
		return lipgloss.NewStyle().Foreground(Gray300).Render(text)
	}
}

// CenterContent centers content in a given width
func CenterContent(content string, width int) string {
	lines := strings.Split(content, "\n")
	var centered []string
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			padding := (width - lineWidth) / 2
			centered = append(centered, strings.Repeat(" ", padding)+line)
		} else {
			centered = append(centered, line)
		}
	}
	return strings.Join(centered, "\n")
}

// CreateCard creates a styled card with header
func CreateCard(title string, content string, width int, selected bool) string {
	var style lipgloss.Style
	if selected {
		style = CardSelectedStyle.Width(width)
	} else {
		style = CardStyle.Width(width)
	}

	header := CardHeaderStyle.Render(" " + title + " ")
	return style.Render(header + "\n" + content)
}

// RenderProgressBar creates an enhanced progress bar
func RenderProgressBar(current, total float64, width int) string {
	if total == 0 {
		return MutedStyle.Render("No data")
	}

	pct := current / total
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}

	filled := int(pct * float64(width))
	empty := width - filled

	// Gradient-like effect using different characters
	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)

	bar := lipgloss.NewStyle().Foreground(Primary).Render(filledBar) +
		lipgloss.NewStyle().Foreground(Gray700).Render(emptyBar)

	percentage := fmt.Sprintf(" %.1f%%", pct*100)

	return bar + lipgloss.NewStyle().Foreground(Gray400).Render(percentage)
}

// RenderHelp renders a help key with description
func RenderHelp(key, description string) string {
	return HelpKeyStyle.Render(key) + " " + HelpDescStyle.Render(description)
}

// RenderStat creates a stat display
func RenderStat(label string, value string, width int) string {
	labelRender := StatLabelStyle.Width(width / 2).Render(label)
	valueRender := StatNumberStyle.Width(width / 2).Align(lipgloss.Right).Render(value)
	return labelRender + valueRender
}

// Backward compatibility aliases
func FormatAmount(amount float64, currency string) string {
	return FormatMoney(amount, currency)
}

func FormatAmountPlain(amount float64, currency string) string {
	return FormatMoneyPlain(amount, currency)
}

func ProgressBar(current, total float64, width int) string {
	return RenderProgressBar(current, total, width)
}
