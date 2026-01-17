package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"echohub-community/internal/models"
)

// NewPost renders the page to create a new post
// @Summary      Show new post page
// @Description  Render the page for creating a new forum post
// @Tags         Posts
// @Produce      html
// @Security     CookieAuth
// @Success      200  {string}  string  "New post page HTML"
// @Failure      401  {object}  models.Error
// @Router       /newPost [get]
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

// AddNewPost creates a new forum post
// @Summary      Create post
// @Description  Add a new post with title, content, and categories
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        newPost  body   models.PostData  true  "Post data"
// @Success      200  {string}  string  "Post created successfully"
// @Failure      400  {object}  models.Response
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /addNewPost [post]
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
	response = webForum.Posts.CheckCategoryIfExist(newPost.Categories)
	if len(response.Messages) != 0 {
		if err := encodeJsonData(w, http.StatusBadRequest, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	newPost.Title = strings.TrimSpace(newPost.Title)
	newPost.Content = strings.TrimSpace(newPost.Content)

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

// UpdatePost renders the page to edit an existing post
// @Summary      Show update post page
// @Description  Render the page for editing an existing forum post
// @Tags         Posts
// @Produce      html
// @Security     CookieAuth
// @Param        ID  query   int  true  "Post ID"
// @Success      200  {string}  string  "Update post page HTML"
// @Failure      400  {object}  models.Error
// @Failure      401  {object}  models.Error
// @Router       /updatePost [get]
func (webForum *WebApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	postId, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		models.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad Request",
			SubMessage: "Invalid PostID",
		}.RenderError(w)
		return
	}

	post, postErr := webForum.Posts.GetPost(user.ID, postId)
	if postErr.Type != "" {
		postErr.RenderError(w)
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

// UpdatingPost updates an existing forum post
// @Summary      Update post
// @Description  Apply changes to title, content, or categories of a post
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        ID  query   int               true  "Post ID"
// @Param        toUpdate  body   models.PostData  true  "Updated post data"
// @Success      200  {string}  string  "Post updated successfully"
// @Failure      400  {object}  models.Response
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /updatingPost [post]
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
	response = webForum.Posts.CheckCategoryIfExist(toUpdate.Categories)
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

	toUpdate.Title = strings.TrimSpace(toUpdate.Title)
	toUpdate.Content = strings.TrimSpace(toUpdate.Content)

	err = webForum.Posts.EditPost(postID, toUpdate.Title, toUpdate.Content, toUpdate.Categories)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, ""); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// DeletePost removes a post from the forum
// @Summary      Delete post
// @Description  Remove an existing post by its ID
// @Tags         Posts
// @Security     CookieAuth
// @Param        ID  query   int  true  "Post ID"
// @Success      302  {string}  string  "Redirect to /"
// @Failure      400  {object}  models.Error
// @Failure      401  {object}  models.Error
// @Failure      500  {object}  models.Error
// @Router       /deletePost [get]
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
