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
	var newPost models.PostData

	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}

	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}

	post, err := webForum.Posts.UpdatePost(user.ID, postId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// func (webForum *WebApp) UpdatingPost(w http.ResponseWriter, r *http.Request) {
// 	var postData PostUpdate
// 	err := json.NewDecoder(r.Body).Decode(&postData)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 		return
// 	}

// 	if postData.Title == "" || len(postData.Title) > 70 || postData.Content == "" || len(postData.Content) > 5000 || len(postData.Categories) == 0 || len(postData.Categories) > 3 {
// 		http.Redirect(w, r, "/update/post?ID="+postData.Id, http.StatusSeeOther)
// 		return
// 	}

// 	id, err := strconv.Atoi(postData.Id)
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusInternalServerError)
// 		return
// 	}

// 	err = webForum.Post.EditPost(id, postData.Title, postData.Content, postData.Categories)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// func (webForum *WebApp) DeletePost(w http.ResponseWriter, r *http.Request) {
// 	user := &models.User{}
// 	var err error

// 	userID, ok := r.Context().Value(models.UserIDKey).(int)
// 	if !ok {
// 		models.Error{
// 			User:       user,
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "Internal Server Error",
// 			SubMessage: "Unable to retrieve user information.",
// 		}.RenderError(w)
// 		return
// 	}

// 	if userID != 0 {
// 		user, err = webForum.Users.FindUserByID(userID)
// 		if err != nil {
// 			if err.Error() == "user not found" {
// 				user = &models.User{
// 					UserType: "guest",
// 				}
// 			} else {
// 				models.Error{
// 					User:       user,
// 					StatusCode: http.StatusInternalServerError,
// 					Message:    "Internal Server Error",
// 					SubMessage: err.Error(),
// 				}.RenderError(w)
// 				return
// 			}
// 		}
// 	} else {
// 		user = &models.User{
// 			UserType: "guest",
// 		}
// 	}
// 	userType, ok := r.Context().Value(userTypeKey).(string)
// 	if !ok {
// 		models.Error{
// 			User:       user,
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "Internal Server Error",
// 			SubMessage: "Unable to retrieve user type.",
// 		}.RenderError(w)
// 		return
// 	}
// 	user.UserType = userType

// 	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusInternalServerError)
// 		return
// 	}

// 	err = webForum.Post.DeletePost(userID, id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }
