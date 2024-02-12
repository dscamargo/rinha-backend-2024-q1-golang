package main

import (
	"context"
	"fmt"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/config"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
)

func main() {
	appConfig := config.Load()

	poolConfig, err := pgxpool.ParseConfig(appConfig.DatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	app := chi.NewRouter()
	//app.Use(middleware.Logger)

	routes.Load(app, dbPool)
	
	log.Fatalln(http.ListenAndServe(":"+appConfig.Port, app))

}
