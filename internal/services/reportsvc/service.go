package reportsvc

import (
    "time"

    "github.com/adelinapasculescu1/atad/internal/repository"
)

type MonthlyReport struct {
    Year           int
    Month          time.Month
    TotalSpent     float64
    ByCategoryName map[string]float64
}

type Service struct {
    txRepo repository.TransactionRepository
}

func NewService(txRepo repository.TransactionRepository) *Service {
    return &Service{txRepo: txRepo}
}

func (s *Service) Monthly(year int, month time.Month) (*MonthlyReport, error) {
    start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
    end := start.AddDate(0, 1, 0)

    total, err := s.txRepo.SumExpenses(start, end)
    if err != nil {
        return nil, err
    }

    byCat, err := s.txRepo.SumExpensesByCategoryName(start, end)
    if err != nil {
        return nil, err
    }

    return &MonthlyReport{
        Year:           year,
        Month:          month,
        TotalSpent:     total,
        ByCategoryName: byCat,
    }, nil
}
