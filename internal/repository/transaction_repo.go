package repository

import (
    "database/sql"
    "time"

    "github.com/adelinapasculescu1/atad/internal/models"
)

type TransactionRepository interface {
    Create(tx *models.Transaction) error
    SumExpensesByCategory(start, end time.Time) (map[int64]float64, error)
    //to be continued
}

type SQLiteTransactionRepository struct {
    db *sql.DB
}

func NewSQLiteTransactionRepository(db *sql.DB) *SQLiteTransactionRepository {
    return &SQLiteTransactionRepository{db: db}
}

func (r *SQLiteTransactionRepository) Create(t *models.Transaction) error {
    query := `
INSERT INTO transactions (date, description, amount, type, category_id, source)
VALUES (?, ?, ?, ?, ?, ?);
`
    var categoryID interface{}
    if t.CategoryID != nil {
        categoryID = *t.CategoryID
    } else {
        categoryID = nil
    }

    res, err := r.db.Exec(
        query,
        t.Date.Format(time.RFC3339),
        t.Description,
        t.Amount,
        string(t.Type),
        categoryID,
        t.Source,
    )
    if err != nil {
        return err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return err
    }

    t.ID = id
    return nil
}

func (r *SQLiteTransactionRepository) SumExpensesByCategory(start, end time.Time) (map[int64]float64, error) {
    query := `
SELECT category_id, COALESCE(SUM(ABS(amount)), 0)
FROM transactions
WHERE type = 'expense'
  AND date >= ?
  AND date < ?
  AND category_id IS NOT NULL
GROUP BY category_id;
`

    rows, err := r.db.Query(query, start.Format(time.RFC3339), end.Format(time.RFC3339))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := map[int64]float64{}
    for rows.Next() {
        var categoryID int64
        var sum float64
        if err := rows.Scan(&categoryID, &sum); err != nil {
            return nil, err
        }
        out[categoryID] = sum
    }
    return out, rows.Err()
}

