package handlers

import (
	"net/http"
	"os"

	"echohub-community/internal/models"

	httpSwagger "github.com/swaggo/http-swagger"
)

type WebApp struct {
	Users         *models.UserModel
	Sessions      *models.SessionModel
	Posts         *models.PostModel
	Comments      *models.CommentModel
	LikesDislikes *models.LikesDislikesModel
}

func (webForum *WebApp) Router() http.Handler {
	mux := http.NewServeMux()

	// Serve "assets" directory
	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))

	// Serve Swagger documentation
	// Swagger JSON endpoint
	mux.HandleFunc("GET /swagger.json", func(w http.ResponseWriter, r *http.Request) {
		swaggerPath := "./docs/swagger.json"
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error":"Swagger documentation not found. Please run 'make docs' to generate it."}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, swaggerPath)
	})

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// Redirect /docs to /swagger/
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	mux.HandleFunc("GET /docs/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

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
	mux.Handle("POST /updateProfile", webForum.AuthMiddleware(http.HandlerFunc(webForum.UpdateProfile)))

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
	mux.Handle("DELETE /deletePost", webForum.AuthMiddleware(http.HandlerFunc(webForum.DeletePost)))

	// Comments route
	mux.Handle("POST /comments", webForum.AuthMiddleware(http.HandlerFunc(webForum.GetComments)))
	mux.Handle("POST /createComment", webForum.AuthMiddleware(http.HandlerFunc(webForum.CreateComment)))

	//Likes & Dislikes Routes
	mux.Handle("POST /like-dislike", webForum.AuthMiddleware(http.HandlerFunc(webForum.LikeDislikeHandler)))

	return mux
}
