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
	multiplexer := http.NewServeMux()
	multiplexer.HandleFunc("GET /", webForum.HomePage)
	multiplexer.HandleFunc("GET /register", webForum.RegisterPage)
	multiplexer.HandleFunc("POST /register", webForum.UserRegister)
	multiplexer.HandleFunc("GET /login", webForum.LoginPage)
	multiplexer.HandleFunc("POST /login", webForum.UserLogin)
	// multiplexer.HandleFunc("GET /get_cookie", webForum.GetCookie)
	multiplexer.HandleFunc("POST /set_cookie", webForum.SetCookie)

	return multiplexer
}
