package services

import (
	"github.com/gorilla/mux"
	"net/http"
)

func MiddlewareSetContentTypeJson() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
