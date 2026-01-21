# DebtQ - Personal Money Tracker

A terminal-based personal finance tracker built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Track expenses, manage debts, monitor investments, and set savings goals - all from your terminal with Obsidian integration.

## Features

### Expense Tracking
- Add and delete expenses with categories
- Categories: food, transport, shopping, utilities, health, entertainment, education, other
- View monthly expense summaries
- Track spending patterns over time

### Borrowing & Lending
- Track money borrowed from others
- Track money lent to others
- **Smart grouping**: Combines all transactions with the same person
- **Net balance calculation**: Shows who owes whom and how much
- **Partial settlements**: Settle specific amounts instead of full transactions
- Due date tracking

### My Net Worth
- Track various investment types:
  - Stocks
  - Mutual Funds
  - Gold & Silver
  - Fixed Deposits
  - PPF
  - Crypto
  - Real Estate
  - Other investments
- Update current values
- Track gains/losses and return percentages

### Savings Goals
- Set savings targets for products you want to buy
- Track progress with visual progress bars
- Add contributions towards goals
- Shows:
  - Days remaining until target date
  - Required monthly savings to reach goal
  - Completion percentage

### Obsidian Integration
- Sync all data to your Obsidian vault as markdown files
- Creates 5 summarized files:
  - `Dashboard.md` - Main overview with links
  - `Expenses.md` - All expenses grouped by month
  - `Debts.md` - All debts grouped by person
  - `NetWorth.md` - Investments grouped by type
  - `Savings.md` - All savings goals

### Stats Dashboard
- Overview of all financial data
- Net worth summary
- Debt position
- Monthly and total expenses
- Savings progress

## Installation

### Prerequisites
- [Go](https://go.dev/dl/) 1.21 or later

### Using Make
```bash
git clone https://github.com/yourusername/debtq.git
cd debtq

# Build and install to ~/.local/bin
make install
```

### Using Install Script
```bash
git clone https://github.com/yourusername/debtq.git
cd debtq

# Run install script
./install.sh
```

### Manual Build
```bash
git clone https://github.com/yourusername/debtq.git
cd debtq

# Build
go build -o debtq ./cmd/main.go

# Move to PATH
mv debtq ~/.local/bin/
```

### Add to PATH
Make sure `~/.local/bin` is in your PATH. Add to your `~/.bashrc` or `~/.zshrc`:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Usage

Start the application:
```bash
debtq
```

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select / Confirm |
| `Esc` | Go back |
| `q` | Quit (from main menu) |

### Expenses View
| Key | Action |
|-----|--------|
| `a` | Add new expense |
| `d` | Delete selected expense |

### Debts View
| Key | Action |
|-----|--------|
| `a` | Add new debt transaction |
| `s` | Settle amount for selected person |

### Net Worth View
| Key | Action |
|-----|--------|
| `a` | Add new investment |
| `u` | Update value of selected investment |
| `d` | Delete selected investment |

### Savings View
| Key | Action |
|-----|--------|
| `a` | Add new savings goal |
| `c` | Add contribution to selected goal |
| `d` | Delete selected goal |

### Form Navigation
| Key | Action |
|-----|--------|
| `Tab` / `↓` | Next field |
| `Shift+Tab` / `↑` | Previous field |
| `Enter` | Save |
| `Esc` | Cancel |

## Configuration

Configuration is stored at `~/.config/debtq/config.json`:

```json
{
  "obsidian_vault_path": "/Users/username/Documents/obsidian-notes/debtq",
  "data_file": "/Users/username/.config/debtq/data.json",
  "currency": "INR"
}
```

### Options
| Option | Description | Default |
|--------|-------------|---------|
| `obsidian_vault_path` | Path to Obsidian vault for markdown export | `~/Documents/obsidian-notes/debtq` |
| `data_file` | Path to JSON data file | `~/.config/debtq/data.json` |
| `currency` | Currency symbol for display | `INR` |

## Data Storage

All data is stored locally in JSON format at `~/.config/debtq/data.json`. The data includes:
- Expenses
- Debt transactions
- Investments
- Savings targets
- Savings contributions

## Make Commands

```bash
make build      # Build the application
make install    # Build and install to ~/.local/bin
make clean      # Remove build artifacts
make uninstall  # Remove from ~/.local/bin
make run        # Build and run the application
make help       # Show help message
```

## Project Structure

```
debtq/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── models/
│   │   └── models.go        # Data models
│   ├── storage/
│   │   ├── storage.go       # JSON data persistence
│   │   └── obsidian.go      # Obsidian markdown generation
│   └── tui/
│       ├── app.go           # Bubble Tea TUI
│       └── styles.go        # Lipgloss styles
├── Makefile
├── install.sh
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [UUID](https://github.com/google/uuid) - Unique ID generation

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
