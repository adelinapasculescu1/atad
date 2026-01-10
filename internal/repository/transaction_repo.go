package repository

import (
    "database/sql"
    "time"
    "strings"

    "github.com/adelinapasculescu1/atad/internal/models"
)

type TransactionRepository interface {
    Create(tx *models.Transaction) error
    SumExpensesByCategory(start, end time.Time) (map[int64]float64, error)
    SumExpenses(start, end time.Time) (float64, error)
    SumExpensesByCategoryName(start, end time.Time) (map[string]float64, error)
    List(filter TransactionFilter) ([]models.Transaction, error)
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

func (r *SQLiteTransactionRepository) SumExpenses(start, end time.Time) (float64, error) {
    query := `
SELECT COALESCE(SUM(ABS(amount)), 0)
FROM transactions
WHERE type = 'expense'
  AND date >= ?
  AND date < ?;
`
    row := r.db.QueryRow(query, start.Format(time.RFC3339), end.Format(time.RFC3339))

    var sum float64
    if err := row.Scan(&sum); err != nil {
        return 0, err
    }
    return sum, nil
}

func (r *SQLiteTransactionRepository) SumExpensesByCategoryName(start, end time.Time) (map[string]float64, error) {
    query := `
SELECT
  COALESCE(c.name, 'Uncategorized') AS category,
  COALESCE(SUM(ABS(t.amount)), 0)
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE t.type = 'expense'
  AND t.date >= ?
  AND t.date < ?
GROUP BY c.name
ORDER BY SUM(ABS(t.amount)) DESC;
`

    rows, err := r.db.Query(query, start.Format(time.RFC3339), end.Format(time.RFC3339))
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[string]float64)
    for rows.Next() {
        var category string
        var sum float64
        if err := rows.Scan(&category, &sum); err != nil {
            return nil, err
        }
        result[category] = sum
    }
    return result, rows.Err()
}

func (r *SQLiteTransactionRepository) List(filter TransactionFilter) ([]models.Transaction, error) {
    query := `
SELECT
  t.id,
  t.date,
  t.description,
  t.amount,
  t.type,
  t.category_id,
  t.source
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE 1=1
`
    args := []any{}

    if filter.From != nil {
        query += " AND t.date >= ?"
        args = append(args, filter.From.Format(time.RFC3339))
    }
    if filter.To != nil {
        query += " AND t.date < ?"
        args = append(args, filter.To.Format(time.RFC3339))
    }
    if filter.Type != nil {
        query += " AND t.type = ?"
        args = append(args, *filter.Type)
    }
    if filter.CategoryName != nil {
        query += " AND c.name = ?"
        args = append(args, *filter.CategoryName)
    }
    if filter.Query != nil {
        query += " AND LOWER(t.description) LIKE ?"
        args = append(args, "%"+strings.ToLower(*filter.Query)+"%")
    }
    if filter.MinAmount != nil {
        query += " AND ABS(t.amount) >= ?"
        args = append(args, *filter.MinAmount)
    }
    if filter.MaxAmount != nil {
        query += " AND ABS(t.amount) <= ?"
        args = append(args, *filter.MaxAmount)
    }

    query += " ORDER BY t.date DESC, t.id DESC"

    limit := filter.Limit
    if limit <= 0 {
        limit = 50
    }
    query += " LIMIT ?"
    args = append(args, limit)

    if filter.Offset > 0 {
        query += " OFFSET ?"
        args = append(args, filter.Offset)
    }

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []models.Transaction
    for rows.Next() {
        var t models.Transaction
        var dateStr string
        var categoryID sql.NullInt64

        if err := rows.Scan(&t.ID, &dateStr, &t.Description, &t.Amount, &t.Type, &categoryID, &t.Source); err != nil {
            return nil, err
        }

        parsedDate, err := time.Parse(time.RFC3339, dateStr)
        if err != nil {
            return nil, err
        }
        t.Date = parsedDate

        if categoryID.Valid {
            cid := categoryID.Int64
            t.CategoryID = &cid
        } else {
            t.CategoryID = nil
        }

        out = append(out, t)
    }

    return out, rows.Err()
}
