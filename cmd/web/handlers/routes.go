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
	forum.HandleFunc("GET /", WebForum.HandleHomePage)
	return forum
}
