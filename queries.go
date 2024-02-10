package main

const (
	GetLastTransactionsQuery = "SELECT t.valor, t.tipo, t.descricao, t.realizado_em AS realizado_em FROM transacoes t WHERE t.cliente_id = $1 ORDER BY t.realizado_em DESC LIMIT 10"
	GetClientBalanceQuery    = "SELECT saldo, limite FROM clientes c WHERE c.id = $1"
)
