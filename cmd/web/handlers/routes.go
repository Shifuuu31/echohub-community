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
	forum.Handle("/assets/", http.StripPrefix("/assets/", fileServer))
	forum.HandleFunc("/", WebForum.HomePage)
	forum.HandleFunc("/post", WebForum.CreatePostPage)
	forum.HandleFunc("/post/create", WebForum.NewPostCreation)
	return forum
}
