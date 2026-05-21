package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router){
			r.Post("/signup/individual", api.handleSignupIndividual)
			r.Post("/signup/company", api.handleSignupCompany)
			r.Post("/login", api.handleLogin)
			r.Post("/logout", api.handleLogout)
		})
	})

}
