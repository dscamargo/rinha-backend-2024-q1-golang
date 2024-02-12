package entities

import "time"

type Statement struct {
	ClientBalance    ClientBalance     `json:"saldo"`
	LastTransactions []LastTransaction `json:"ultimas_transacoes"`
}

type ClientBalance struct {
	Total         int       `json:"total"`
	Limit         int       `json:"limite"`
	StatementDate time.Time `json:"data_extrato"`
}

type LastTransaction struct {
	Amount      int             `json:"valor"`
	Kind        TransactionKind `json:"tipo"`
	Description string          `json:"descricao"`
	DueAt       time.Time       `json:"realizado_em"`
}
