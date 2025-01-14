package middleware

import (
	"context"
	"net/http"
	"forum/internal/models"
)

func AuthMiddleware(sessionModel *models.SessionModel, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := sessionModel.ValidateSession(sessionToken.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}