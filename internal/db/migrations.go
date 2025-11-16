package db

import "database/sql"

func Migrate(db *sql.DB) error {
    //to be continued
    query := `
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    type TEXT NOT NULL, -- "income" sau "expense"
    category_id INTEGER,
    source TEXT NOT NULL
);
`
    _, err := db.Exec(query)
    return err
}
