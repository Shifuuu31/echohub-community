package web

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Users *models.UserModel
	Sessions *models.SessionModel
}

func (webForum *WebApp) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", webForum.HomePage)
	mux.HandleFunc("GET /register", webForum.RegisterPage)
	mux.HandleFunc("POST /register", webForum.UserRegister)
	mux.HandleFunc("GET /login", webForum.LoginPage)
	mux.HandleFunc("POST /login", webForum.UserLogin)
	mux.HandleFunc("GET /get_cookie", webForum.GetCookie)
	mux.HandleFunc("POST /set_cookie", webForum.SetCookie)

	return mux
}
