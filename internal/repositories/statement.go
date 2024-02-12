package repositories

import (
	"context"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	GetLastTransactionsQuery = "SELECT t.valor, t.tipo, t.descricao, t.realizado_em AS realizado_em FROM transacoes t WHERE t.cliente_id = $1 ORDER BY t.realizado_em DESC LIMIT 10"
	GetClientBalanceQuery    = "SELECT saldo, limite FROM clientes c WHERE c.id = $1"
)

type StatementRepository struct {
	dbPool *pgxpool.Pool
}

func (r *StatementRepository) List(ctx context.Context, clientId int) (entities.Statement, error) {
	var balanceOutput entities.ClientBalance
	err := r.dbPool.QueryRow(ctx, GetClientBalanceQuery, clientId).Scan(&balanceOutput.Total, &balanceOutput.Limit)
	if err != nil {
		return entities.Statement{}, err
	}

	lastTransactions := make([]entities.LastTransaction, 0)
	lastTransactionsRows, err := r.dbPool.Query(ctx, GetLastTransactionsQuery, clientId)
	if err != nil {
		return entities.Statement{}, err
	}

	for lastTransactionsRows.Next() {
		var lastTransaction entities.LastTransaction
		err := lastTransactionsRows.Scan(&lastTransaction.Amount, &lastTransaction.Kind, &lastTransaction.Description, &lastTransaction.DueAt)
		if err != nil {
			return entities.Statement{}, err
		}
		lastTransactions = append(lastTransactions, lastTransaction)
	}

	return entities.Statement{
		ClientBalance: entities.ClientBalance{
			Total:         balanceOutput.Total,
			Limit:         balanceOutput.Limit,
			StatementDate: time.Now().UTC(),
		},
		LastTransactions: lastTransactions,
	}, nil
}

func NewStatementRepository(dbPool *pgxpool.Pool) *StatementRepository {
	return &StatementRepository{dbPool}
}
