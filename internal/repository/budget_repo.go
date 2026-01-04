package repository

import (
    "database/sql"

    "github.com/adelinapasculescu1/atad/internal/models"
)

type BudgetRepository interface {
    Upsert(b models.Budget) (*models.Budget, error)
    List() ([]models.Budget, error)
    GetByCategoryID(categoryID int64) (*models.Budget, error)
}

type SQLiteBudgetRepository struct {
    db *sql.DB
}

func NewSQLiteBudgetRepository(db *sql.DB) *SQLiteBudgetRepository {
    return &SQLiteBudgetRepository{db: db}
}

func (r *SQLiteBudgetRepository) Upsert(b models.Budget) (*models.Budget, error) {
    query := `
INSERT INTO budgets (category_id, amount, period)
VALUES (?, ?, ?)
ON CONFLICT(category_id) DO UPDATE SET
    amount = excluded.amount,
    period = excluded.period;
`
    _, err := r.db.Exec(query, b.CategoryID, b.Amount, b.Period)
    if err != nil {
        return nil, err
    }

    return r.GetByCategoryID(b.CategoryID)
}

func (r *SQLiteBudgetRepository) GetByCategoryID(categoryID int64) (*models.Budget, error) {
    row := r.db.QueryRow(`SELECT id, category_id, amount, period FROM budgets WHERE category_id = ?;`, categoryID)
    var b models.Budget
    if err := row.Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Period); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &b, nil
}

func (r *SQLiteBudgetRepository) List() ([]models.Budget, error) {
    rows, err := r.db.Query(`SELECT id, category_id, amount, period FROM budgets ORDER BY id ASC;`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []models.Budget
    for rows.Next() {
        var b models.Budget
        if err := rows.Scan(&b.ID, &b.CategoryID, &b.Amount, &b.Period); err != nil {
            return nil, err
        }
        out = append(out, b)
    }
    return out, rows.Err()
}
