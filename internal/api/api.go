package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Router *chi.Mux
	Sessions *scs.SessionManager
}