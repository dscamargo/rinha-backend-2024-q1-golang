package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/apperrors"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/entities"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

type InsertTransactionHandler struct {
	svc *services.InsertTransactionService
}

func (h *InsertTransactionHandler) Execute(w http.ResponseWriter, r *http.Request) {
	clientIdStr := chi.URLParam(r, "id")
	clientId, err := strconv.Atoi(clientIdStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if clientId > 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var body entities.Transaction
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := body.Validate(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	body.ClientId = clientId
	output, err := h.svc.Execute(ctx, body)

	if err != nil {
		if errors.Is(err, apperrors.ErrBalanceIsNull) || errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		log.Println("erro ao atualizar saldo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, output)
}

func NewInsertTransactionHandler(svc *services.InsertTransactionService) *InsertTransactionHandler {
	return &InsertTransactionHandler{svc}
}
