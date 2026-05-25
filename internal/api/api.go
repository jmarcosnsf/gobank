package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jmarcosnsf/gobank/internal/services"
)

type Api struct {
	Router *chi.Mux
	Sessions *scs.SessionManager
	HolderService services.HolderService
	AccountService services.AccountService
}