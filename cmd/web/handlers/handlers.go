package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"forum/internal/models"
)

type FetchCredentials struct {
	Post_id  int    `json:"postID"`
	Category string `json:"category"`
}

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

// handler to get categories
func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}

	categories, err := webForum.Post.GetCategories()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! Cannot retrieve categories at the moment."}.RenderError(w)
		return
	}

	if err := Template.ExecuteTemplate(w, "home.html", categories); err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! Failed to render the home page template."}.RenderError(w)
	}
}

// Handler to get maxID
func (webForm *WebApp) GetMaxID(w http.ResponseWriter, r *http.Request) {
	maxID, err := webForm.Post.GetMaxID()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! Unable to fetch the maximum ID at the moment."}.RenderError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(maxID); err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! Failed to encode the response data."}.RenderError(w)
		return
	}
}

// handler to get posts
func (webForm *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	var credentials FetchCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if credentials.Category == "" {
		return
	}

	post, err := webForm.Post.GetPosts(credentials.Post_id, credentials.Category)
	if err != nil || post.PostId == 0 {
		fmt.Fprintf(w, "null")
		return
	}

	var postsBuffer bytes.Buffer
	encoder := json.NewEncoder(&postsBuffer)
	if err := encoder.Encode(post); err != nil {
		http.Error(w, "Failed to encode header object", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, postsBuffer.String())
}

// desplay create-post page
func (webForum *WebApp) CreatePost(w http.ResponseWriter, r *http.Request) {
	Categories, err := webForum.Post.GetCategories()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! " + err.Error()}.RenderError(w)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-creation.html", Categories); err != nil {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}
}

func (webForm *WebApp) PostCreation(w http.ResponseWriter, r *http.Request) {
	New := models.Post{
		PostTitle:   r.FormValue("title"),
		PostContent: r.FormValue("content"),
	}
	categoriesForm := r.Form["categoryElement"]

	if New.PostTitle == "" || len(New.PostTitle) > 70 || New.PostContent == "" || len(New.PostContent) > 5000 || len(categoriesForm) == 0 || len(categoriesForm) > 3 {
		http.Redirect(w, r, "/createPost", http.StatusSeeOther)
		return
	}

	ids, err := webForm.Post.GetIdsCategories(categoriesForm)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! " + err.Error()}.RenderError(w)
		return
	}

	idPost, err := webForm.Post.CreatePost(New.PostTitle, New.PostContent)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! " + err.Error()}.RenderError(w)
		return
	}

	err = webForm.Post.AddcategoryPost(idPost, ids)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error", SubMessage: "Oops! " + err.Error()}.RenderError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func (webForum *WebApp) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/post/delete" {
// 		w.WriteHeader(http.StatusNotFound)
// 		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
// 			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
// 			return
// 		}
// 		return
// 	}
// 	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusInternalServerError)
// 		return
// 	}

// 	// fmt.Printf("%V",id)
// 	err = webForum.Post.DeletePost(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

func (WebForum *WebApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postUpdate" {
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

func (WebApp *WebApp) PostUpdate(w http.ResponseWriter, r *http.Request) {
	New := models.Post{
		PostTitle:   r.FormValue("title"),
		PostContent: r.FormValue("content"),
	}

	categoriesForm := r.Form["categories[]"]

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
