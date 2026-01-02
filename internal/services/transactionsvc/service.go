package transactionsvc

import (
    "fmt"
    "time"

    "github.com/adelinapasculescu1/atad/internal/models"
    "github.com/adelinapasculescu1/atad/internal/repository"
)

type Service struct {
    txRepo repository.TransactionRepository
}

func NewService(txRepo repository.TransactionRepository) *Service {
    return &Service{txRepo: txRepo}
}

type AddManualInput struct {
    Date        time.Time
    Description string
    Amount      float64
    Type        models.TransactionType
    CategoryID  *int64 // to be modified
}

func (s *Service) AddManual(input AddManualInput) error {
    if input.Description == "" {
        return fmt.Errorf("description is required")
    }
    if input.Amount == 0 {
        return fmt.Errorf("amount cannot be 0")
    }
    if input.Type != models.TransactionTypeIncome && input.Type != models.TransactionTypeExpense {
        return fmt.Errorf("invalid transaction type: %s", input.Type)
    }

    tx := &models.Transaction{
        Date:        input.Date,
        Description: input.Description,
        Amount:      input.Amount,
        Type:        input.Type,
        CategoryID:  input.CategoryID,
        Source:      "manual",
    }

    if err := s.txRepo.Create(tx); err != nil {
        return fmt.Errorf("saving transaction: %w", err)
    }

    return nil
}
