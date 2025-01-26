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
	forum.HandleFunc("GET /", WebForum.HomePage)
	forum.HandleFunc("GET /max-id", WebForum.GetMaxID)
	forum.HandleFunc("POST /post", WebForum.GetPosts)
	forum.HandleFunc("GET /createPost", WebForum.CreatePost)
	forum.HandleFunc("POST /createPost", WebForum.PostCreation)
	forum.HandleFunc("GET /postUpdate", WebForum.UpdatePost)
	forum.HandleFunc("POST /postUpdate", WebForum.UpdatePost)
	
	// forum.HandleFunc("GET /post/update", WebForum.UpdatePostPageHandler)
	// forum.HandleFunc("POST /post/update/edit", WebForum.PostUpdateHandler)

	return forum
}
