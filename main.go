package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"strconv"
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

	app := chi.NewRouter()
	//app.Use(middleware.Logger)

	app.Post("/clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		postTransactionsController(w, r, dbConn)
	})
	app.Get("/clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		getClientBalanceController(w, r, dbConn)
	})
	log.Fatalln(http.ListenAndServe(":"+config.port, app))

}

func postTransactionsController(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Pool) {
	clientIdStr := chi.URLParam(r, "id")
	clientId, err := strconv.Atoi(clientIdStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if clientId > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var body transactionInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if (body.Tipo != "c" && body.Tipo != "d") || body.Valor == 0 || len(body.Descricao) < 1 || len(body.Descricao) > 10 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
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
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		log.Println("erro ao atualizar saldo: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if updatedBalance == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	render.JSON(w, r, map[string]int{
		"saldo":  *updatedBalance,
		"limite": *updatedLimit,
	})
}

func getClientBalanceController(w http.ResponseWriter, r *http.Request, dbConn *pgxpool.Pool) {
	clientIdStr := chi.URLParam(r, "id")
	clientId, err := strconv.Atoi(clientIdStr)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if clientId > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var balanceOutput balanceResponse
	err = dbConn.QueryRow(ctx, GetClientBalanceQuery, clientId).Scan(&balanceOutput.Total, &balanceOutput.Limite)
	if err != nil {
		log.Println("erro ao buscar saldo do cliente: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lastTransactions := make([]lastTransactionsResponse, 0)
	lastTransactionsRows, err := dbConn.Query(ctx, GetLastTransactionsQuery, clientId)
	if err != nil {
		log.Println("erro ao buscar ultimas transacoes: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for lastTransactionsRows.Next() {
		var lastTransaction lastTransactionsResponse
		err := lastTransactionsRows.Scan(&lastTransaction.Valor, &lastTransaction.Tipo, &lastTransaction.Descricao, &lastTransaction.RealizadoEm)
		if err != nil {
			log.Println("extrato Scan: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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

	render.JSON(w, r, response)
}
