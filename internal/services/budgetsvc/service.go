package budgetsvc

import (
    "fmt"
    "time"

    "github.com/adelinapasculescu1/atad/internal/models"
    "github.com/adelinapasculescu1/atad/internal/repository"
)

const warningThreshold = 0.90

type Status string

const (
    StatusOK       Status = "OK"
    StatusWarning  Status = "WARNING"
    StatusExceeded Status = "EXCEEDED"
)

type BudgetStatus struct {
    CategoryName string
    BudgetAmount float64
    SpentAmount  float64
    Remaining    float64
    Status       Status
}

type Service struct {
    budgetsRepo  repository.BudgetRepository
    categoryRepo repository.CategoryRepository
    txRepo       repository.TransactionRepository
}

func NewService(
    budgetsRepo repository.BudgetRepository,
    categoryRepo repository.CategoryRepository,
    txRepo repository.TransactionRepository,
) *Service {
    return &Service{
        budgetsRepo:  budgetsRepo,
        categoryRepo: categoryRepo,
        txRepo:       txRepo,
    }
}

func (s *Service) SetBudget(categoryName string, amount float64, period string) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be > 0")
    }
    if period == "" {
        period = "monthly"
    }
    if period != "monthly" {
        return fmt.Errorf("unsupported period: %s (only monthly supported for now)", period)
    }

    cat, err := s.categoryRepo.GetByName(categoryName)
    if err != nil {
        return err
    }
    if cat == nil {
        return fmt.Errorf("category not found: %s", categoryName)
    }

    _, err = s.budgetsRepo.Upsert(models.Budget{
        CategoryID: cat.ID,
        Amount:     amount,
        Period:     period,
    })
    return err
}

func (s *Service) GetMonthlyStatus(year int, month time.Month) ([]BudgetStatus, error) {
    budgets, err := s.budgetsRepo.List()
    if err != nil {
        return nil, err
    }

    start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
    end := start.AddDate(0, 1, 0)

    spentByCatID, err := s.txRepo.SumExpensesByCategory(start, end)
    if err != nil {
        return nil, err
    }

    out := make([]BudgetStatus, 0, len(budgets))
    for _, b := range budgets {
        cat, err := s.categoryRepo.GetByID(b.CategoryID)
        if err != nil {
            return nil, err
        }
        name := "(unknown)"
        if cat != nil {
            name = cat.Name
        }

        spent := spentByCatID[b.CategoryID]
        remaining := b.Amount - spent

        st := StatusOK
        ratio := 0.0
        if b.Amount > 0 {
            ratio = spent / b.Amount
        }
        if ratio >= 1.0 {
            st = StatusExceeded
        } else if ratio >= warningThreshold {
            st = StatusWarning
        }

        out = append(out, BudgetStatus{
            CategoryName: name,
            BudgetAmount: b.Amount,
            SpentAmount:  spent,
            Remaining:    remaining,
            Status:       st,
        })
    }

    return out, nil
}
