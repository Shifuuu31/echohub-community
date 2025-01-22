package handlers

import (
	"net/http"

	"forum/internal/models"
)

type WebApp struct {
	Post *models.PostModel
}

func (WebForum *WebApp) Routes() http.Handler {
	forum := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./assets"))
	forum.Handle("GET /assets/", http.StripPrefix("/assets/", fileServer))
	forum.HandleFunc("GET /", WebForum.HomePageHandler)
	forum.HandleFunc("GET /post", WebForum.GetPostsHandler)
	forum.HandleFunc("GET /create-post", WebForum.CreatePostPageHandler)
	forum.HandleFunc("POST /create", WebForum.NewPostCreationHandler)
	// forum.HandleFunc("GET /post/update", WebForum.UpdatePostPageHandler)
	// forum.HandleFunc("POST /post/update/edit", WebForum.PostUpdateHandler)

	return forum
}
