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
	forum.HandleFunc("POST /createPost", WebForum.Creation)
	forum.HandleFunc("GET /postUpdate", WebForum.Update)
	forum.HandleFunc("POST /postUpdate", WebForum.Updating)
	forum.HandleFunc("POST /deletePost", WebForum.Delete)

	return forum
}
