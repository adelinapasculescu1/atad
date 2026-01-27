**Command-line personal finance manager using Go**

command-line application written in Go that tracks personal income and expenses directly from the terminal.  
It supports importing transactions from bank statements, manual entry, categorization, budgeting, reporting and searching.

**Features**
- Import transactions from **CSV** and **OFX** files
- Manually add income and expense transactions
- Manage categories
- Set monthly budgets per category and get alerts
- Generate monthly spending reports
- Search and filter transactions
- SQLite-based local storage

**Clone and run**

```bash
git clone https://github.com/adelinapasculescu1/atad.git
cd atad
go run ./cmd/main.go --help
```

**CLI Commands**

*Import Transactions*
atad import --file transactions.csv --format csv
atad import --file bank.ofx --format ofx

*Add Manual Transaction*
atad add \
  --type expense \
  --amount 50 \
  --description "Groceries" \
  --date 2025-01-15 \
  --category "Food"

Flags:
--type income | expense
--amount transaction amount
--description transaction description
--date YYYY-MM-DD (optional)
--category category name (optional)

*Categories operations*
atad category add --name "Groceries"
atad category list

*Set budgets*
atad budget set --category "Groceries" --amount 300

*Check budget status for a specific month*
atad budget status --month 2025-01

*Reports*
atad report monthly --month 2025-01

*Search and filter*
atad search \
  --from 2025-01-01 \
  --to 2025-01-31 \
  --type expense \
  --category "Groceries" \
  --q "lidl" \
  --min 10 \
  --max 100
  
Available filters:
--from, --to –> date range
--type –>income or expense
--category -> category name
--q –> keyword search in description
--min, --max -> amount range
--limit -> max results
