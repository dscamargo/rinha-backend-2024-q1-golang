package main

const (
	GetLastTransactionsQuery = "SELECT t.valor, t.tipo, t.descricao, t.created_at AS realizado_em FROM transacoes t WHERE t.cliente_id = $1 ORDER BY t.created_at DESC LIMIT 10"
	GetClientBalanceQuery    = "SELECT saldo, limite FROM clientes c WHERE c.id = $1"
)
