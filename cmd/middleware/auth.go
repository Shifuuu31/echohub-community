package middleware

import (
	"context"
	"net/http"

	"forum/internal/models"
)

func AuthMiddleware(sessionModel *models.SessionModel, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("userSession")
		if err != nil {
			models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot Get Cookie"}.RenderError(w)
			return
		}

		userID, err := sessionModel.ValidateSession(sessionToken.Value)
		if err != nil {
			models.Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized", SubMessage: "Invalid Session"}.RenderError(w)
			return
		}

		ctx := context.WithValue(r.Context(), "UserID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
