package web

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
}

func (webForum *WebApp) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))

	mux.Handle("GET /", webForum.AuthMiddleware(http.HandlerFunc(webForum.HomePage)))

	mux.HandleFunc("GET /register", webForum.RegisterPage)
	mux.HandleFunc("POST /confirmRegister", webForum.UserRegister)
	mux.Handle("GET /login", webForum.AuthMiddleware(http.HandlerFunc(webForum.LoginPage)))
	// mux.HandleFunc("POST /login", webForum.UserLogin)
	mux.HandleFunc("POST /confirmLogin", webForum.ConfirmLogin)
	mux.HandleFunc("GET /logout", webForum.UserLogout)
	// mux.HandleFunc("POST /set_cookie", webForum.SetCookie)

	return mux
}
