package routes

import (
	"github.com/dscamargo/rinha-2024-q1-golang/internal/handlers"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/repositories"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Load(app *chi.Mux, dbPool *pgxpool.Pool) {
	statementRepository := repositories.NewStatementRepository(dbPool)
	transactionRepository := repositories.NewTransactionRepository(dbPool)

	getStatementService := services.NewGetStatementService(statementRepository)
	insertTransactionService := services.NewInsertTransactionService(transactionRepository)

	getStatementHandler := handlers.NewGetStatementHandler(getStatementService)
	insertTransactionHandler := handlers.NewInsertTransactionHandler(insertTransactionService)

	app.Post("/clientes/{id}/transacoes", insertTransactionHandler.Execute)
	app.Get("/clientes/{id}/extrato", getStatementHandler.Execute)
}
