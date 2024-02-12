package services

import (
	"context"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/entities"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/repositories"
)

type InsertTransactionService struct {
	repo *repositories.TransactionRepository
}

func (svc *InsertTransactionService) Execute(ctx context.Context, body entities.Transaction) (entities.InsertTransactionOutput, error) {
	return svc.repo.Insert(ctx, body)
}

func NewInsertTransactionService(repo *repositories.TransactionRepository) *InsertTransactionService {
	return &InsertTransactionService{repo}
}
