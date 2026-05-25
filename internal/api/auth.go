package api

import (
	"net/http"

	"github.com/jmarcosnsf/gobank/internal/jsonutils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "holder_id") {
			jsonutils.EncondeJson(w, r, http.StatusUnauthorized, map[string]any{"error": "must be logged in"})
			return
		}
		next.ServeHTTP(w, r)
	})
}