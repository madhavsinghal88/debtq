package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/debtq/debtq/internal/config"
	"github.com/debtq/debtq/internal/storage"
	"github.com/debtq/debtq/internal/tui"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Ensure Obsidian directory exists
	if err := cfg.EnsureObsidianDir(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Obsidian directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize storage
	store, err := storage.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	// Create and run TUI
	model := tui.New(cfg, store)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
