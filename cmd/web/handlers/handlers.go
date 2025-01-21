package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"forum/internal/models"
)

var (
	Template       = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))
	categoriesForm = []string{}
)

func (webForum *WebApp) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	categories, err := webForum.Post.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "home.html", categories); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	Categories, err := webForum.Post.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-create.html", Categories); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) NewPostCreationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	New := models.Post{
		PostTitle:   r.FormValue("title"),
		PostContent: r.FormValue("content"),
	}
	categoriesForm = r.Form["categories[]"]

	ids, err := webForm.Post.GetIdsCategories(categoriesForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idPost, err := webForm.Post.CreatePost(New.PostTitle, New.PostContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	err = webForm.Post.AddcategoryPost(idPost, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (webForum *WebApp) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	// fmt.Printf("%V",id)
	err = webForum.Post.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (WebForum *WebApp) UpdatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/update" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	title, content, selected_categorys, err := WebForum.Post.UpdatePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Categorys, err := WebForum.Post.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		ID                  int
		Title               string
		Content             string
		Categories          []models.Category
		Categories_selected []string
	}{
		ID:                  id,
		Title:               title,
		Content:             content,
		Categories_selected: selected_categorys,
		Categories:          Categorys,
	}

	if err := Template.ExecuteTemplate(w, "post-update.html", data); err != nil {
		http.Error(w, "Error loading UpdatePage"+err.Error(), http.StatusInternalServerError)
	}
}

func (WebApp *WebApp) PostUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/update/edit" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	New := models.Post{
		PostTitle:   r.FormValue("title"),
		PostContent: r.FormValue("content"),
	}

	categoriesForm = r.Form["categories[]"]

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = WebApp.Post.EditPost(id, New.PostTitle, New.PostContent, categoriesForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
