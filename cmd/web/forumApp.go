package web

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
	LikesDislikes *models.LikesDislikesModel 
}

func (webForum *WebApp) Routes() http.Handler {
	mux := http.NewServeMux()

	// Serve "assets" directory
	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))

	// authentication middleware
	mux.Handle("GET /", webForum.AuthMiddleware(http.HandlerFunc(webForum.HomePage)))

	// Registration routes
	mux.Handle("GET /register", webForum.AuthMiddleware(http.HandlerFunc(webForum.RegisterPage)))
	mux.HandleFunc("POST /confirmRegister", webForum.UserRegister)

	// Login routes
	mux.Handle("GET /login", webForum.AuthMiddleware(http.HandlerFunc(webForum.LoginPage)))
	mux.HandleFunc("POST /confirmLogin", webForum.ConfirmLogin)

	// Logout route
	mux.HandleFunc("GET /logout", webForum.UserLogout)

	// Like/Dislike route
	mux.Handle("POST /like-dislike", webForum.AuthMiddleware(http.HandlerFunc(webForum.LikeDislikeHandler)))
	mux.HandleFunc("GET /likes-dislikes", webForum.GetLikesDislikesHandler)

	return mux
}
