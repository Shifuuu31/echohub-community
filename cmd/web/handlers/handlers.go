package handlers

import (
	"net/http"
	"text/template"
)

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	posts, err := webForum.Post.GetPosts()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	}
	if err := Template.ExecuteTemplate(w, "home2.html", posts); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/create-post" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	// if r.Method != http.MethodPost {
	// 	http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	// }
	if err := Template.ExecuteTemplate(w, "create_post.html", nil); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}
