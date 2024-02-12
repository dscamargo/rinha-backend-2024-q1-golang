package handlers

import (
	"context"
	"github.com/dscamargo/rinha-2024-q1-golang/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"time"
)

type GetStatementHandler struct {
	svc *services.GetStatementService
}

func (h *GetStatementHandler) Execute(w http.ResponseWriter, r *http.Request) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := h.svc.Execute(ctx, clientId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, output)
}

func NewGetStatementHandler(svc *services.GetStatementService) *GetStatementHandler {
	return &GetStatementHandler{svc}
}
