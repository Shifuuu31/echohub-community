package handlers

import (
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (webForum *WebApp) NewPost(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType != "authenticated" {
		userErr = models.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			SubMessage: "Please try to login",
		}

		userErr.RenderError(w)
		return
	}

	categories, catsErr := webForum.Posts.GetCategories()
	if catsErr.Type == "server" {
		catsErr.RenderError(w)
		return
	}

	CreatePostData := struct {
		User       *models.User
		Categories []models.Category
	}{
		User:       user,
		Categories: categories,
	}
	models.RenderPage(w, "newPost.html", CreatePostData)
}

func (webForum *WebApp) AddNewPost(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}

	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var newPost models.PostData

	if decodeErr := decodeJsonData(r, &newPost); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	response := models.CheckNewPost(newPost)
	if len(response.Messages) != 0 {
		if err := encodeJsonData(w, http.StatusBadRequest, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	postID, err := webForum.Posts.CreatePost(user.ID, newPost.Title, newPost.Content)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = webForum.Posts.AddCategoriesPost(postID, newPost.Categories)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, ""); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType != "authenticated" {
		userErr = models.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			SubMessage: "Please try to login",
		}

		userErr.RenderError(w)
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		models.Error{
			StatusCode: http.StatusBadRequest,
			Message: "Bad Request",
			SubMessage: "Invalid PostID",
		}.RenderError(w)
		return
	}

	post, err := webForum.Posts.GetPost(user.ID, postId)
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal Server Error",
		}.RenderError(w)
		return
	}

	categories, catsErr := webForum.Posts.GetCategories()
	if catsErr.Type == "server" {
		catsErr.RenderError(w)
		return
	}

	data := struct {
		User       *models.User
		Post_info  models.Post
		Categories []models.Category
	}{
		User:       user,
		Post_info:  post,
		Categories: categories,
	}

	models.RenderPage(w, "updatePost.html", data)
}

func (webForum *WebApp) UpdatingPost(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}

	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var toUpdate models.PostData

	if decodeErr := decodeJsonData(r, &toUpdate); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	response := models.CheckNewPost(toUpdate)
	if len(response.Messages) != 0 {
		if err := encodeJsonData(w, http.StatusBadRequest, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}

	err = webForum.Posts.EditPost(postID, toUpdate.Title, toUpdate.Content, toUpdate.Categories)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, ""); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) DeletePost(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType != "authenticated" {
		models.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			SubMessage: "Please try to login",
		}.RenderError(w)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		models.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			SubMessage: "Invalid PostID",
		}.RenderError(w)
		return
	}

	err = webForum.Posts.DeletePost(user.ID, postID)
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}.RenderError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
