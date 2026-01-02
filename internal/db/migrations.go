package db

import "database/sql"

func Migrate(db *sql.DB) error {
    queries := []string{
        `
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    type TEXT NOT NULL,
    category_id INTEGER,
    source TEXT NOT NULL
);
`,
        `
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);
`,
        `
CREATE INDEX IF NOT EXISTS idx_transactions_category_id
ON transactions(category_id);
`,
    }

    for _, q := range queries {
        if _, err := db.Exec(q); err != nil {
            return err
        }
    }

    return nil
}
