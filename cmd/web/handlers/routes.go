package handlers

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
	Post     *models.PostModel
}

func (webForum *WebApp) Router() http.Handler {
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

	mux.HandleFunc("GET /maxId", webForum.GetMaxID)
	mux.HandleFunc("POST /post", webForum.GetPosts)

	mux.Handle("GET /createPost", webForum.AuthMiddleware(http.HandlerFunc(webForum.CreatePost)))
	mux.Handle("POST /createPost", webForum.AuthMiddleware(http.HandlerFunc(webForum.Creation)))

	mux.Handle("GET /updatePost", webForum.AuthMiddleware(http.HandlerFunc(webForum.UpdatePost)))
	mux.HandleFunc("POST /updatePost", webForum.Updating)

	mux.Handle("GET /deletePost", webForum.AuthMiddleware(http.HandlerFunc(webForum.DeletePost)))
	// mux.HandleFunc("GET /deletePost", webForum.DeletePost)

	mux.HandleFunc("POST /comments", webForum.GetComments)

	return mux
}
