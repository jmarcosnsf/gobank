package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/jmarcosnsf/gobank/internal/jsonutils"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("GOBANK_CSRF_KEY")),
		csrf.Secure(false),
		csrf.Path("/"),
	)
	env := os.Getenv("GOBANK_ENV")

	if env != "dev" {
		api.Router.Use(csrfMiddleware)
	}

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrftoken", api.HandleGetCSRFToken)

			r.Post("/signup/individual", api.handleSignupIndividual)
			r.Post("/signup/company", api.handleSignupCompany)
			r.Post("/login", api.handleLogin)

			r.Group(func(r chi.Router) {
				r.Use(api.AuthMiddleware)

				r.Post("/logout", api.handleLogout)

				r.Route("/account", func(r chi.Router) {
					r.Post("/", api.handleCreateAccount)
					r.Get("/{id}/balance", api.handleGetBalance)
					r.Post("/{id}/deposit", api.handleDeposit)
					r.Post("/{id}/withdrawal", api.handleWithdraw)
					r.Post("/transfer", api.handleTransfer)
					r.Delete("/{id}", api.handleCloseAccount)
				})
			})
		})
	})
}

func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	jsonutils.EncondeJson(w, r, http.StatusOK, map[string]any{"csrf_token": token})
}
