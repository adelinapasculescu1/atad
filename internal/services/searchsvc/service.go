package searchsvc

import (
    "github.com/adelinapasculescu1/atad/internal/models"
    "github.com/adelinapasculescu1/atad/internal/repository"
)

type Service struct {
    txRepo repository.TransactionRepository
}

func NewService(txRepo repository.TransactionRepository) *Service {
    return &Service{txRepo: txRepo}
}

func (s *Service) Search(filter repository.TransactionFilter) ([]models.Transaction, error) {
    return s.txRepo.List(filter)
}
