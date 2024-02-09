package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

type balanceResponse struct {
	Total       int       `json:"total"`
	Limite      int       `json:"limite"`
	DataExtrato time.Time `json:"data_extrato"`
}

type lastTransactionsResponse struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadoEm time.Time `json:"realizado_em"`
}

type balanceResponseOutput struct {
	Saldo             balanceResponse            `json:"saldo"`
	UltimasTransacoes []lastTransactionsResponse `json:"ultimas_transacoes"`
}

type transactionInput struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type appConfig struct {
	port        string
	databaseUrl string
}

func getEnvOrDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func main() {
	config := appConfig{
		port:        getEnvOrDefault("PORT", "8080"),
		databaseUrl: getEnvOrDefault("DATABASE_URL", "postgresql://pg:pg@localhost:5432/rinha"),
	}

	poolConfig, err := pgxpool.ParseConfig(config.databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbConn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{})

	app.Post("/clientes/:id/transacoes", func(ctx fiber.Ctx) error {
		return postTransactionsController(ctx, dbConn)
	})
	app.Get("/clientes/:id/extrato", func(ctx fiber.Ctx) error {
		return getClientBalanceController(ctx, dbConn)
	})

	log.Fatalln(app.Listen(":" + config.port))
}

func postTransactionsController(c fiber.Ctx, dbConn *pgxpool.Pool) error {
	clientId, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if clientId > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	var body transactionInput
	if err := c.Bind().JSON(&body); err != nil {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	if (body.Tipo != "c" && body.Tipo != "d") || body.Valor == 0 || len(body.Descricao) < 1 || len(body.Descricao) > 10 {
		return c.SendStatus(http.StatusUnprocessableEntity)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	valor := body.Valor
	if body.Tipo == "d" {
		valor = valor * -1
	}

	var updatedBalance *int
	var updatedLimit *int
	err = dbConn.QueryRow(ctx, "call criar_transacao($1, $2, $3, $4)", clientId, valor, body.Tipo, body.Descricao).Scan(&updatedBalance, &updatedLimit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(http.StatusUnprocessableEntity).Send([]byte{})
		}
		log.Println("erro ao atualizar saldo: ", err)
		return c.Status(http.StatusInternalServerError).Send([]byte{})
	}

	if updatedBalance == nil {
		return c.Status(http.StatusUnprocessableEntity).Send([]byte{})
	}

	return c.Status(http.StatusOK).JSON(map[string]int{
		"saldo":  *updatedBalance,
		"limite": *updatedLimit,
	})
}

func getClientBalanceController(c fiber.Ctx, dbConn *pgxpool.Pool) error {
	clientId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte{})
	}

	if clientId > 5 {
		return c.SendStatus(http.StatusNotFound)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var balanceOutput balanceResponse
	err = dbConn.QueryRow(ctx, GetClientBalanceQuery, clientId).Scan(&balanceOutput.Total, &balanceOutput.Limite)
	if err != nil {
		log.Println("erro ao buscar saldo do cliente: ", err)
		return c.Status(http.StatusInternalServerError).Send([]byte{})
	}

	lastTransactions := make([]lastTransactionsResponse, 0)
	lastTransactionsRows, err := dbConn.Query(ctx, GetLastTransactionsQuery, clientId)
	if err != nil {
		log.Println("erro ao buscar ultimas transacoes: ", err)
		return c.Status(http.StatusInternalServerError).Send([]byte{})
	}

	for lastTransactionsRows.Next() {
		var lastTransaction lastTransactionsResponse
		err := lastTransactionsRows.Scan(&lastTransaction.Valor, &lastTransaction.Tipo, &lastTransaction.Descricao, &lastTransaction.RealizadoEm)
		if err != nil {
			log.Println("extrato Scan: ", err)
			return c.Status(http.StatusInternalServerError).Send([]byte{})
		}
		lastTransactions = append(lastTransactions, lastTransaction)
	}

	response := balanceResponseOutput{
		Saldo: balanceResponse{
			Total:       balanceOutput.Total,
			Limite:      balanceOutput.Limite,
			DataExtrato: time.Now(),
		},
		UltimasTransacoes: lastTransactions,
	}

	return c.Status(http.StatusOK).JSON(response)
}
