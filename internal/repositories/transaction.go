package repositories

import (
	"context"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/apperrors"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	dbPool *pgxpool.Pool
}

func (r *TransactionRepository) Insert(ctx context.Context, body entities.Transaction) (entities.InsertTransactionOutput, error) {
	valor := body.Amount
	if body.Kind == entities.TransactionKindDebit {
		valor = valor * -1
	}

	var updatedBalance, updatedLimit *int
	err := r.dbPool.QueryRow(ctx, "CALL criar_transacao($1, $2, $3, $4)", body.ClientId, valor, body.Kind, body.Description).Scan(&updatedBalance, &updatedLimit)
	if err != nil {
		return entities.InsertTransactionOutput{}, err
	}

	if updatedBalance == nil {
		return entities.InsertTransactionOutput{}, apperrors.ErrBalanceIsNull
	}

	return entities.InsertTransactionOutput{
		Balance: *updatedBalance,
		Limit:   *updatedLimit,
	}, nil
}

func NewTransactionRepository(dbPool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{dbPool}
}
