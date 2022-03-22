package middleware

import (
	"net/http"
)

func TypeApplicationJSON(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
    })
}