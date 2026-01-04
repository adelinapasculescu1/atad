package repository

import (
    "database/sql"

    "github.com/adelinapasculescu1/atad/internal/models"
)

type CategoryRepository interface {
    Create(name string) (*models.Category, error)
    GetByName(name string) (*models.Category, error)
    GetByID(id int64) (*models.Category, error)
    List() ([]models.Category, error)
}

type SQLiteCategoryRepository struct {
    db *sql.DB
}

func NewSQLiteCategoryRepository(db *sql.DB) *SQLiteCategoryRepository {
    return &SQLiteCategoryRepository{db: db}
}

func (r *SQLiteCategoryRepository) Create(name string) (*models.Category, error) {
    res, err := r.db.Exec(`INSERT INTO categories (name) VALUES (?);`, name)
    if err != nil {
        return nil, err
    }
    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }
    return &models.Category{ID: id, Name: name}, nil
}

func (r *SQLiteCategoryRepository) GetByName(name string) (*models.Category, error) {
    row := r.db.QueryRow(`SELECT id, name FROM categories WHERE name = ?;`, name)

    var c models.Category
    if err := row.Scan(&c.ID, &c.Name); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &c, nil
}

func (r *SQLiteCategoryRepository) GetByID(id int64) (*models.Category, error) {
    row := r.db.QueryRow(`SELECT id, name FROM categories WHERE id = ?;`, id)

    var c models.Category
    if err := row.Scan(&c.ID, &c.Name); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return &c, nil
}


func (r *SQLiteCategoryRepository) List() ([]models.Category, error) {
    rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY name ASC;`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []models.Category
    for rows.Next() {
        var c models.Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    return out, rows.Err()
}
