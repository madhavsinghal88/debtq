# DebtQ - Personal Money Tracker

A terminal-based personal finance tracker built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Track expenses, manage debts, monitor investments, and set savings goals - all from your terminal with Obsidian integration.

## Features

### Expense Tracking
- Add and delete expenses with categories
- Categories: food, transport, shopping, utilities, health, entertainment, education, other
- View monthly expense summaries
- Track spending patterns over time
- Math expressions supported in amount fields (e.g., "100+50")

### Borrowing & Lending
- Track money borrowed from others
- Track money lent to others
- **Smart grouping**: Combines all transactions with the same person
- **Net balance calculation**: Shows who owes whom and how much
- **Transaction selection**: Choose specific transactions to settle
- **Partial settlements**: Settle specific amounts instead of full transactions
- **Settlement notes**: Add descriptions to remember why/how you settled (e.g., "Cash payment", "UPI transfer")
- **Settlement history**: View all settled transactions with notes and amounts
- **Person history**: Detailed view of all transactions with a specific person
- **Custom transaction dates**: Enter the actual date when money was borrowed/lent

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
- Debt position (borrowed vs lent)
- Monthly and total expenses
- Savings progress tracking

## Installation

### Prerequisites
- [Go](https://go.dev/dl/) 1.25 or later
- Alternatively, download the pre-built binary from [Releases](https://github.com/debtq/debtq/releases)

### Using Make
```bash
git clone https://github.com/madhavsinghal88/debtq.git
cd debtq

# Build and install to ~/.local/bin
make install
```

### Using Install Script
```bash
git clone https://github.com/madhavsinghal88/debtq.git
cd debtq

# Run install script
./install.sh
```

### Manual Build
```bash
git clone https://github.com/madhavsinghal88/debtq.git
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
| `s` | Settle amount for selected person (select transaction first) |
| `h` | View settlement history (all settled transactions) |
| `i` | View person history (detailed view for selected person) |

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

## Settlement Workflow

When settling debts, you now have more control:

1. **Select Person** - In Debts view, select a person and press `s`
2. **Choose Transaction** - Select the specific transaction you want to settle
3. **Enter Amount** - Enter the settlement amount (defaults to full amount)
4. **Add Note** - Add a settlement note (e.g., "Paid via UPI", "Cash payment")
5. **Confirm** - Press Enter to complete

**Partial Settlements**: If you settle less than the full amount, the transaction is automatically split:
- Original transaction keeps the remaining amount (still active)
- New settled transaction is created for the settled portion

## Viewing History

### Settlement History (`h`)
Shows all settled transactions across all people:
- Settlement date
- Original and settled amounts
- Settlement notes
- Transaction descriptions

### Person History (`i`)
Detailed view for a selected person showing:
- **Summary**: Total lent/borrowed, settled amounts, net balance
- **Active Transactions**: All unsettled debts with descriptions
- **Settlement History**: Chronological list of all settlements with notes

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
- Debt transactions (with settlement tracking)
- Investments
- Savings targets
- Savings contributions

### Debt Transaction Format

```json
{
  "id": "abc123",
  "type": "lent",
  "person_name": "John",
  "amount": 1000.00,
  "description": "Lunch money",
  "date": "2026-01-15T00:00:00Z",
  "is_settled": true,
  "settled_date": "2026-02-09T10:30:00Z",
  "settlement_amount": 500.00,
  "settlement_note": "Paid via UPI"
}
```

**Settlement Fields:**
- `is_settled`: Whether the transaction is settled
- `settled_date`: When it was settled
- `settlement_amount`: Actual amount settled (may differ from original for partial settlements)
- `settlement_note`: Description of why/how it was settled

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
