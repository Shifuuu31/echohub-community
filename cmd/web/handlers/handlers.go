package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
)

type contextKey string

var (
	userIDKey   contextKey = "UserID"
	userTypeKey contextKey = "UserType"
)

func (webForum *WebApp) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := "guest"
		var userID int

		sessionToken, err := r.Cookie("userSession")
		if err == nil {
			userID, err = webForum.Sessions.ValidateSession(sessionToken.Value)
			if err == nil {
				userType = "authenticated"
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, userID)
		ctx = context.WithValue(ctx, userTypeKey, userType)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	var err error

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	user.UserType = userType
	if r.URL.Path != "/" {
		err := models.Error{
			User:       user,
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! The page you are looking for does not exist",
		}
		err.RenderError(w)
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
		User       *models.User
		Categories []models.Category
	}{
		User:       user,
		Categories: categories,
	}

	models.RenderPage(w, "home.html", homeData)
}

func (webForum *WebApp) LoginPage(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	var err error

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	if userType == "authenticated" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	models.RenderPage(w, "login.html", nil)
}

type UserCredentials struct {
	UserName   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe"`
}

func (webForum *WebApp) ConfirmLogin(w http.ResponseWriter, r *http.Request) {
	var credentials UserCredentials
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userID, errors := webForum.Users.ValidateUserCredentials(credentials.UserName, credentials.Password)
	fmt.Println(userID)
	if userID > 0 {
		newSession, err := webForum.Sessions.GenerateNewSession(userID, credentials.RememberMe)
		if err != nil {
			http.Error(w, "Internal Server Error1", http.StatusInternalServerError)
			return
		}

		newCookie, err := webForum.Sessions.InsertOrUpdateSession(newSession)
		if err != nil {
			http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &newCookie)

		sendJsontoHeader(w, []string{"Login successful!"})
	} else {
		sendJsontoHeader(w, errors)
	}
}

func (webForum *WebApp) UserLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("userSession")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err = webForum.Sessions.DeleteSession(sessionCookie.Value)
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Failed to delete session",
		}.RenderError(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "userSession",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (webForum *WebApp) RegisterPage(w http.ResponseWriter, r *http.Request) {
	models.RenderPage(w, "register.html", nil)
}

type NewUserInfo struct {
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RepeatedPass string `json:"rPassword"`
}

func (webForum *WebApp) UserRegister(w http.ResponseWriter, r *http.Request) {
	var newUserinfo NewUserInfo
	err := json.NewDecoder(r.Body).Decode(&newUserinfo)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	newUser, errors := webForum.Users.ValidateNewUser(newUserinfo.UserName, newUserinfo.Email, newUserinfo.Password, newUserinfo.RepeatedPass)
	if len(errors) == 0 {
		if err := webForum.Users.InsertUser(newUser); err != nil {
			models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot insert new user"}.RenderError(w)
			return
		}
		sendJsontoHeader(w, []string{"User Registred successfully!"})
	} else {
		sendJsontoHeader(w, errors)
	}
	fmt.Println("error", errors)
}

func (webForum *WebApp) GetMaxID(w http.ResponseWriter, r *http.Request) {
	maxID, err := webForum.Post.GetMaxID()
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

	// sendJsontoHeader(w, maxID)
	if err := json.NewEncoder(w).Encode(maxID); err != nil { // to be replaced look on top
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! Failed to fetch the maximum ID at the moment.",
		}.RenderError(w)
		return
	}
}

type FetchPosts struct {
	Post_id  int    `json:"postID"`
	Category string `json:"category"`
}

func (webForum *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	var postData FetchPosts
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Oops! Failed to get posts", http.StatusInternalServerError)
		return
	}
	if postData.Category == "" {
		return
	}

	post, err := webForum.Post.GetPosts(postData.Post_id, postData.Category)
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

func (webForum *WebApp) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	// var err error
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

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	user.UserType = userType

	CreatePostData := struct {
		User       *models.User
		Categories []models.Category
	}{
		User:       user,
		Categories: categories,
	}
	models.RenderPage(w, "post-creation.html", CreatePostData)
}

func (webForum *WebApp) Creation(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	var err error

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	user.UserType = userType

	var postData PostUpdate
	err = json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if postData.Title == "" || len(postData.Title) > 70 || postData.Content == "" || len(postData.Content) > 5000 || len(postData.Categories) == 0 || len(postData.Categories) > 3 {
		http.Redirect(w, r, "/createPost", http.StatusSeeOther)
		return
	}

	ids, err := webForum.Post.GetIdsCategories(postData.Categories)
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! " + err.Error(),
		}.RenderError(w)
		return
	}

	idPost, err := webForum.Post.CreatePost(userID, postData.Title, postData.Content)
	if err != nil {
		models.Error{
			User:       &models.User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "500 Internal Server Error",
			SubMessage: "Oops! " + err.Error(),
		}.RenderError(w)
		return
	}

	err = webForum.Post.AddCategoriesPost(idPost, ids)
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

func (webForum *WebApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	var err error

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	user.UserType = userType

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	post, err := webForum.Post.UpdatePost(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Categorys, err := webForum.Post.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User       *models.User
		Post_info  models.Post
		Categories []models.Category
	}{
		User:       user,
		Post_info:  post,
		Categories: Categorys,
	}
	fmt.Println("",post.PostCategories,)

	models.RenderPage(w, "post-update.html", data)
}

type PostUpdate struct {
	Id         string   `json:"id`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

func (webForum *WebApp) Updating(w http.ResponseWriter, r *http.Request) {
	var postData PostUpdate
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if postData.Title == "" || len(postData.Title) > 70 || postData.Content == "" || len(postData.Content) > 5000 || len(postData.Categories) == 0 || len(postData.Categories) > 3 {
		http.Redirect(w, r, "/update/post?ID="+postData.Id, http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(postData.Id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = webForum.Post.EditPost(id, postData.Title, postData.Content, postData.Categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (webForum *WebApp) DeletePost(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	var err error

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user information.",
		}.RenderError(w)
		return
	}

	if userID != 0 {
		user, err = webForum.Users.FindUserByID(userID)
		if err != nil {
			if err.Error() == "user not found" {
				user = &models.User{
					UserType: "guest",
				}
			} else {
				models.Error{
					User:       user,
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					SubMessage: err.Error(),
				}.RenderError(w)
				return
			}
		}
	} else {
		user = &models.User{
			UserType: "guest",
		}
	}
	userType, ok := r.Context().Value(userTypeKey).(string)
	if !ok {
		models.Error{
			User:       user,
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Unable to retrieve user type.",
		}.RenderError(w)
		return
	}
	user.UserType = userType

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = webForum.Post.DeletePost(userID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// tools-------------------------------------

func sendJsontoHeader(w http.ResponseWriter, obj interface{}) error {
	fmt.Println("OBJ:\x1b[1;31m", obj, "\x1b[1;39m")
	w.Header().Set("Content-Type", "application/json")
	var jsonData bytes.Buffer
	encoder := json.NewEncoder(&jsonData)
	if err := encoder.Encode(obj); err != nil {
		return errors.New("failed to encode object: " + err.Error())
	}
	fmt.Fprintf(w, jsonData.String())

	return nil
}
