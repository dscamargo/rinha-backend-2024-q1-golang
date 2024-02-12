package services

import (
	"context"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/entities"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/repositories"
)

type GetStatementService struct {
	repo *repositories.StatementRepository
}

func (svc *GetStatementService) Execute(ctx context.Context, clientId int) (entities.Statement, error) {
	return svc.repo.List(ctx, clientId)
}

func NewGetStatementService(repo *repositories.StatementRepository) *GetStatementService {
	return &GetStatementService{repo}
}
