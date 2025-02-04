package handlers

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Users    *models.UserModel
	Sessions *models.SessionModel
	Posts     *models.PostModel
	Comments *models.CommentModel
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

	// ProfileSettings routes
	mux.Handle("GET /profileSettings", webForum.AuthMiddleware(http.HandlerFunc(webForum.ProfileSettings)))
	// mux.HandleFunc("POST /UpdateProfile", webForum.UpdateProfile)
	
	// MaxID route
	mux.HandleFunc("POST /maxId", webForum.MaxID)
	
	// GetPosts route
	mux.Handle("POST /posts", webForum.AuthMiddleware(http.HandlerFunc(webForum.GetPosts)))
	
	// NewPost routes
	mux.Handle("GET /newPost", webForum.AuthMiddleware(http.HandlerFunc(webForum.NewPost)))
	mux.Handle("POST /addNewPost", webForum.AuthMiddleware(http.HandlerFunc(webForum.AddNewPost)))
	
	// UpdatePost routes
	mux.Handle("GET /updatePost", webForum.AuthMiddleware(http.HandlerFunc(webForum.UpdatePost)))
	mux.Handle("POST /updatingPost", webForum.AuthMiddleware(http.HandlerFunc(webForum.UpdatingPost)))
	
	// DeletePost route
	mux.Handle("GET /deletePost", webForum.AuthMiddleware(http.HandlerFunc(webForum.DeletePost)))
	
	// Comments route
	mux.HandleFunc("POST /comments", webForum.GetComments)
	mux.Handle("POST /createComment", webForum.AuthMiddleware(http.HandlerFunc(webForum.CreateComment)))

	return mux
}
