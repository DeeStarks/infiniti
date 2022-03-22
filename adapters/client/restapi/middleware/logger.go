package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deestarks/infiniti/utils"
)

func Logger(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		handler.ServeHTTP(w, r)

		// Logging
		msg := fmt.Sprintf("%s %s %s. Response Time: %v", r.Method, r.RequestURI, r.Proto, time.Since(start))
		utils.LogMessage(msg)
    })
}