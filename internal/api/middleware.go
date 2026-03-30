package api

import (
	"log/slog"
	"net/http"
)

func LogRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			url    = r.URL.String()
		)

		slog.Info("request", "method", method, "url", url)
		next.ServeHTTP(w, r)
	})
}
