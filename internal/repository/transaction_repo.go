package repository

import (
    "database/sql"
    "time"

    "github.com/adelinapasculescu1/atad/internal/models"
)

type TransactionRepository interface {
    Create(tx *models.Transaction) error
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
