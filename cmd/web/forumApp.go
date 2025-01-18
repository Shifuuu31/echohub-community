package web

import (
	"net/http"

	"forum/internal/models"
	// auth "forum/cmd/middleware"

)

type WebApp struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
}

func (webForum *WebApp) Routes() http.Handler {
	mux := http.NewServeMux()

	// Use mux.Handle for routes that require middleware
	mux.HandleFunc("GET /", webForum.HomePage)
	mux.HandleFunc("GET /register", webForum.RegisterPage)
	mux.HandleFunc("POST /register", webForum.UserRegister)
	mux.HandleFunc("GET /login", webForum.LoginPage)
	mux.HandleFunc("POST /login", webForum.UserLogin)
	mux.HandleFunc("POST /set_cookie", webForum.SetCookie)

	return mux
}
