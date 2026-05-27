package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmarcosnsf/gobank/internal/jsonutils"
	"github.com/jmarcosnsf/gobank/internal/services"
	"github.com/jmarcosnsf/gobank/internal/usecase/account"
	"github.com/shopspring/decimal"
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

func (api *Api) handleDeposit(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(rawID)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid account id"})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[account.DepositReq](r)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	holderID, holderType, err := api.getCurrentHolder(r.Context())
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	err = api.AccountService.Deposit(r.Context(), accountID, holderID, holderType, decimal.NewFromFloat(data.Amount))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			jsonutils.EncondeJson(w, r, http.StatusNotFound, map[string]any{"error": "account not found"})
		case errors.Is(err, services.ErrAccountNotOwned):
			jsonutils.EncondeJson(w, r, http.StatusForbidden, map[string]any{"error": "account does not belong to you"})
		case errors.Is(err, services.ErrAccountNotActive):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "account is not active"})
		default:
			slog.Error("deposit failed", "error", err)
			jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		}
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "deposit successful"})
}

func (api *Api) handleCloseAccount(w http.ResponseWriter, r *http.Request) {
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

	err = api.AccountService.CloseAccount(r.Context(), accountID, holderID, holderType)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			jsonutils.EncondeJson(w, r, http.StatusNotFound, map[string]any{"error": "account not found"})
		case errors.Is(err, services.ErrAccountNotOwned):
			jsonutils.EncondeJson(w, r, http.StatusForbidden, map[string]any{"error": "account does not belong to you"})
		case errors.Is(err, services.ErrAccountNotActive):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "account is not active"})
		case errors.Is(err, services.ErrAccountHasBalance):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "account still has balance"})
		default:
			slog.Error("close account failed", "error", err)
			jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		}
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "account closed"})
}

func (api *Api) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(rawID)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid account id"})
		return
	}

	data, problems, err := jsonutils.DecodeValidJson[account.WithdrawalReq](r)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	holderID, holderType, err := api.getCurrentHolder(r.Context())
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	err = api.AccountService.Withdraw(r.Context(), accountID, holderID, holderType, decimal.NewFromFloat(data.Amount))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			jsonutils.EncondeJson(w, r, http.StatusNotFound, map[string]any{"error": "account not found"})
		case errors.Is(err, services.ErrAccountNotOwned):
			jsonutils.EncondeJson(w, r, http.StatusForbidden, map[string]any{"error": "account does not belong to you"})
		case errors.Is(err, services.ErrAccountNotActive):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "account is not active"})
		case errors.Is(err, services.ErrInsufficientFunds):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "insufficient funds"})
		default:
			slog.Error("withdraw failed", "error", err)
			jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		}
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "withdraw successful"})
}

func (api *Api) handleTransfer(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[account.TransferReq](r)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	fromID, err := uuid.Parse(data.FromAccountID)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid from_account_id"})
		return
	}
	toID, err := uuid.Parse(data.ToAccountID)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid to_account_id"})
		return
	}

	holderID, holderType, err := api.getCurrentHolder(r.Context())
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	err = api.AccountService.Transfer(r.Context(), fromID, toID, holderID, holderType, decimal.NewFromFloat(data.Amount))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAccountNotFound):
			jsonutils.EncondeJson(w, r, http.StatusNotFound, map[string]any{"error": "account not found"})
		case errors.Is(err, services.ErrAccountNotOwned):
			jsonutils.EncondeJson(w, r, http.StatusForbidden, map[string]any{"error": "source account does not belong to you"})
		case errors.Is(err, services.ErrAccountNotActive):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "one of the accounts is not active"})
		case errors.Is(err, services.ErrInsufficientFunds):
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "insufficient funds"})
		case errors.Is(err, services.ErrSameAccountTransfer):
			jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "cannot transfer to the same account"})
		default:
			slog.Error("transfer failed", "error", err)
			jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		}
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "transfer successful"})
}