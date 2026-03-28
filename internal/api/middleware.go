package api

import (
	"context"
	"imdb/internal/security"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "user_id"

func AuthMiddleware(ts *security.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if len(header) < 7 || header[:7] != "Bearer " {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := header[7:]
			userID, err := ts.Verify(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContextOr500(w http.ResponseWriter, r *http.Request) string {
	if userID, ok := r.Context().Value(userIDKey).(string); ok {
		return userID
	}

	http.Error(w, "Internal Error", 500)
	return ""
}
