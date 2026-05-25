package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmarcosnsf/gobank/internal/jsonutils"
	"github.com/jmarcosnsf/gobank/internal/services"
)

func (api *Api) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	holderID, holderType, err := api.getCurrentHolder(r.Context())
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	id, err := api.AccountService.CreateAccount(r.Context(), holderID, holderType)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusCreated, map[string]any{"account_id": id})
}

func (api *Api) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(rawID)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid account id"})
		return
	}

	holderID, holderType, err := api.getCurrentHolder(r.Context())
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	balance, err := api.AccountService.GetBalance(r.Context(), accountID, holderID, holderType)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			jsonutils.EncondeJson(w, r, http.StatusNotFound, map[string]any{"error": "account not found"})
		case errors.Is(err, services.ErrAccountNotOwned):
			jsonutils.EncondeJson(w, r, http.StatusForbidden, map[string]any{"error": "account does not belong to you"})
		default:
			slog.Error("get balance failed", "error", err)
			jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		}
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"balance": balance})
}

func (api *Api) handleDeposit(w http.ResponseWriter, r *http.Request)      {}
func (api *Api) handleCloseAccount(w http.ResponseWriter, r *http.Request) {}
