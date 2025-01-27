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

type FetchPosts struct {
	Post_id  int    `json:"postID"`
	Category string `json:"category"`
}

type PostUpdate struct {
	Id         string   `json:"id`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

// handler to get categories
func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! the page you looking for does not exist",
		}.RenderError(w)
		return
	}

	categories, err := webForum.Post.GetCategories()
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Cannot retrieve categories at the moment",
		}.RenderError(w)
		return
	}

	homeData := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}

	if err := Template.ExecuteTemplate(w, "home.html", homeData); err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Failed to render the home page template.",
		}.RenderError(w)
	}
}

// Handler to get maxID
func (webForm *WebApp) GetMaxID(w http.ResponseWriter, r *http.Request) {
	maxID, err := webForm.Post.GetMaxID()
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Failed to fetch the maximum ID at the moment.",
		}.RenderError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(maxID); err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Failed to fetch the maximum ID at the moment.",
		}.RenderError(w)
		return
	}
}

// handler to get posts
func (webForm *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	var postData FetchPosts
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Oops! Failed to get posts", http.StatusInternalServerError)
		return
	}
	if postData.Category == "" {
		return
	}

	post, err := webForm.Post.GetPosts(postData.Post_id, postData.Category)
	if err != nil || post.PostId == 0 {
		fmt.Fprintf(w, "null")
		return
	}

	var postsBuffer bytes.Buffer
	encoder := json.NewEncoder(&postsBuffer)
	if err := encoder.Encode(post); err != nil {
		http.Error(w, "Oops! Failed to get posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, postsBuffer.String())
}

// create-post page
func (webForum *WebApp) CreatePost(w http.ResponseWriter, r *http.Request) {
	Categories, err := webForum.Post.GetCategories()
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Cannot retrieve categories at the moment",
		}.RenderError(w)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-creation.html", Categories); err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! the page you looking for does not exist",
		}.RenderError(w)
		return
	}
}

func (webForm *WebApp) Creation(w http.ResponseWriter, r *http.Request) {
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
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! " + err.Error(),
		}.RenderError(w)
		return
	}

	idPost, err := webForm.Post.CreatePost(New.PostTitle, New.PostContent)
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! " + err.Error(),
		}.RenderError(w)
		return
	}

	err = webForm.Post.AddCategoriesPost(idPost, ids)
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! " + err.Error(),
		}.RenderError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// update post page
func (WebForum *WebApp) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		// models.Error{
		// 	User:       &models.User{},
		// 	StatusCode: http.StatusInternalServerError,
		// 	Message:    "500 Internal Server Error",
		// 	SubMessage: "Oops! " + err.Error(),
		// }.RenderError(w)
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	post, err := WebForum.Post.UpdatePost(id)
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
		Post_info  models.Post
		Categories []models.Category
	}{
		Post_info:  post,
		Categories: Categorys,
	}

	if err := Template.ExecuteTemplate(w, "post-update.html", data); err != nil {
		http.Error(w, "Error loading UpdatePage", http.StatusInternalServerError)
	}
}

// update post
func (WebApp *WebApp) Updating(w http.ResponseWriter, r *http.Request) {
	var postData PostUpdate
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if postData.Title == "" || postData.Content == "" || len(postData.Categories) == 0 || len(postData.Categories) > 3 {
		return
	}

	id, err := strconv.Atoi(postData.Id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = WebApp.Post.EditPost(id, postData.Title, postData.Content, postData.Categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// delete post
func (webForum *WebApp) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = webForum.Post.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
