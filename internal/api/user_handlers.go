package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/jmarcosnsf/gobank/internal/jsonutils"
	"github.com/jmarcosnsf/gobank/internal/services"
	"github.com/jmarcosnsf/gobank/internal/usecase/holder"
	"github.com/shopspring/decimal"
)

func (api *Api) handleSignup(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[holder.CreateIndividualReq](r)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	dob, err := time.Parse("2006-01-02", data.DateOfBirth)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, map[string]any{"error": "invalid date format"})
		return
	}

	id, err := api.HolderService.CreateIndividualHolder(
		r.Context(),
		data.FullName,
		data.Cpf,
		data.Email,
		data.Phone,
		data.Category,
		dob,
		decimal.NewFromFloat(data.MonthlyIncome),
		data.Password,
	)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmailOrCpf) {
			jsonutils.EncondeJson(w, r, http.StatusConflict, map[string]any{"error": "cpf or email already exists"})
			return
		}
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	jsonutils.EncondeJson(w, r, http.StatusCreated, map[string]any{"holder_id": id})
}

func (api *Api) handleLogin(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[holder.LoginReq](r)
	if err != nil {
		jsonutils.EncondeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, holderType, err := api.HolderService.Authenticate(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncondeJson(w, r, http.StatusBadRequest, map[string]any{"error": "invalid email or password"})
			return
		}

		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		jsonutils.EncondeJson(w, r, http.StatusInternalServerError, map[string]any{"error": "something went wrong"})
		return
	}

	api.Sessions.Put(r.Context(), "holder_id", id.String())
	api.Sessions.Put(r.Context(), "holder_type", string(holderType))

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "succesfully logged in"})
}

func (api *Api) handleLogout(w http.ResponseWriter, r *http.Request) {
	api.Sessions.Destroy(r.Context())

	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"message": "sucessfully logged out"})
}