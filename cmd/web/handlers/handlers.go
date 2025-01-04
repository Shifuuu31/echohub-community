package handlers

import (
	"html/template"
	"net/http"

	"forum/internal/models"
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

	if err := Template.ExecuteTemplate(w, "home2.html", posts); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	categories, err := webForm.Post.GetCategorys()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "create_post.html", categories); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) NewPostCreation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	New := models.Post{
		Title:        r.FormValue("title"),
		Post_content: r.FormValue("content"),
	}
	// categories := r.Form["categorys"]

	// ids, err := webForm.Post.GetIdsCategories(categories)
	// if err != nil {
	// http.Error(w, err.Error(), http.StatusInternalServerError)
	// return
	// }

	err := webForm.Post.CreatePost(New.Title, New.Post_content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
