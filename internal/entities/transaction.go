package entities

import "github.com/dscamargo/rinha-2024-q1-golang/internal/apperrors"

type TransactionKind string

const (
	TransactionKindDebit  TransactionKind = "d"
	TransactionKindCredit TransactionKind = "c"
)

type Transaction struct {
	ClientId    int             `json:"cliente_id"`
	Amount      int             `json:"valor"`
	Kind        TransactionKind `json:"tipo"`
	Description string          `json:"descricao"`
}

func (e *Transaction) Validate() error {
	if (e.Kind != TransactionKindDebit && e.Kind != TransactionKindCredit) || e.Amount == 0 || len(e.Description) < 1 || len(e.Description) > 10 {
		return apperrors.ErrInvalidBody
	}
	return nil
}

type InsertTransactionOutput struct {
	Balance int `json:"saldo"`
	Limit   int `json:"limite"`
}
