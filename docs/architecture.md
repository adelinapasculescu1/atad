# Command-line personal finance manager using Go

## System architecture

[ User ]
   | 
   v
[ CLI / TUI layer ]
   |
   v
[ Services layer ]
   |
   v
[ Repository layer ]
   |
   v
[ SQLite database ]

| layer          | description                                              |
| -------------- | ------------------------------------------------------------- |
| **CLI / TUI**  | User interaction, terminal commands, interactive interface    |
| **Services**   | Business logic: importing, categorizing, budgeting, reporting |
| **Repository** | Data persistence, querying, filtering                         |
| **Database**   | SQLite storage for transactions, budgets, categories, rules   |

## Project structure 

ATAD/
  cmd/
    main.go
  internal/
    cli/
    tui/
    config/
    db/
    models/
    repository/
    services/
      importsvc/
      categsvc/
      budgetsvc/
      reportsvc/
      searchsvc/
    ui/
  tests/
  docs/
    architecture.md
    secisions.md
    images/


## Data models

Transaction
- ID
- date
- description
- amount
- type (income/expense)
- categoryID (nullable)
- source (manual, import)

Category
- ID
- name

CategoryRule
- ID
- categoryID
- pattern (regex)
- field to match (e.g., description)
- priority (to break ties)

Budget
- ID
- categoryID
- amount (limit)
- period (weekly/monthly)