package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

func (webForum *WebApp) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int
		var userType string
		var sessionErr models.Error
		ctx := r.Context()

		userCookie, err := r.Cookie("userSession")
		if err != nil {
			userType = "guest"
		} else {
			userID, sessionErr = webForum.Sessions.ValidateSession(userCookie.Value)
			if sessionErr.Type == "server" {
				sessionErr.RenderError(w)
				return
			}
			if userID > 0 {
				userType = "authenticated"
			}
		}
		ctx = context.WithValue(ctx, models.UserIDKey, userID)
		ctx = context.WithValue(ctx, models.UserTypeKey, userType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if r.URL.Path != "/" {
		models.Error{
			User:       user,
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! The page you are looking for does not exist",
		}.RenderError(w)
		return
	}

	categories, catsErr := webForum.Post.GetCategories()
	if catsErr.Type == "server" {
		catsErr.RenderError(w)
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
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType == "authenticated" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	models.RenderPage(w, "login.html", nil)
}

func (webForum *WebApp) ConfirmLogin(w http.ResponseWriter, r *http.Request) {
	credentials := struct {
		UserName   string `json:"username"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}{}

	if decodeErr := decodeJsonData(r, &credentials); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var userID int
	response := models.Response{}
	userID, response.Messages = webForum.Users.ValidateUserCredentials(credentials.UserName, credentials.Password)
	if userID > 0 {
		newSession, newSessionErr := webForum.Sessions.GenerateNewSession(userID, credentials.RememberMe)
		if newSessionErr.Type == "server" {
			if err := encodeJsonData(w, http.StatusInternalServerError, newSessionErr); err != nil {
				http.Error(w, "failed to encode object.", http.StatusInternalServerError)
			}
			return
		}

		newCookie, newSessionErr := webForum.Sessions.InsertOrUpdateSession(newSession)
		if newSessionErr.Type == "server" {
			if err := encodeJsonData(w, http.StatusInternalServerError, newSessionErr); err != nil {
				http.Error(w, "failed to encode object.", http.StatusInternalServerError)
			}
			return
		}
		http.SetCookie(w, &newCookie)

		response.Messages = append(response.Messages, "Login successful!")
		if err := encodeJsonData(w, http.StatusOK, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	} else {
		if err := encodeJsonData(w, http.StatusUnauthorized, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	}
}

func (webForum *WebApp) UserLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("userSession")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	delSessionErr := webForum.Sessions.DeleteSession(sessionCookie.Value)
	if delSessionErr.Type == "server" {
		delSessionErr.RenderError(w)
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
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType == "authenticated" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	models.RenderPage(w, "register.html", nil)
}

func (webForum *WebApp) UserRegister(w http.ResponseWriter, r *http.Request) {
	var newUserinfo models.NewUserInfo

	if decodeErr := decodeJsonData(r, &newUserinfo); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	newUser, response := webForum.Users.ValidateNewUser(newUserinfo)
	if len(response.Messages) == 0 {
		if err := webForum.Users.InsertUser(newUser); err != nil {
			models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot insert new user"}.RenderError(w)
			return
		}
		response.Messages = append(response.Messages, "User Registred successfully!")
		if err := encodeJsonData(w, http.StatusOK, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	} else {
		if err := encodeJsonData(w, http.StatusBadRequest, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	}
}

func (webForum *WebApp) MaxID(w http.ResponseWriter, r *http.Request) {
	maxID, maxIdErr := webForum.Post.GetMaxId()
	if maxIdErr.Type == "server" {
		maxIdErr.RenderError(w)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, maxID); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	postsData := struct {
		StartId  int    `json:"start"`
		Category string `json:"category"`
	}{}

	if decodeErr := decodeJsonData(r, &postsData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	posts, postErr := webForum.Post.GetPosts(postsData.StartId, postsData.Category)
	if postErr.Type != "" {
		if postErr.Type == "server" {
			postErr.RenderError(w)
			return
		}

		if err := encodeJsonData(w, postErr.StatusCode, postErr); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	if len(posts) == 0 {
		postErr = models.Error{
			StatusCode: http.StatusContinue,
			Message:    "No posts available",
			Type:       "client",
		}

		if err := encodeJsonData(w, postErr.StatusCode, postErr); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return

	}

	if err := encodeJsonData(w, http.StatusOK, posts); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// func (webForum *WebApp) CreatePost(w http.ResponseWriter, r *http.Request) {
// 	user := &models.User{}
// 	// var err error
// 	categories, err := webForum.Post.GetCategories()
// 	if err != nil {
// 		models.Error{
// 			User:       &models.User{},
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "500 Internal Server Error",
// 			SubMessage: "Oops! Cannot retrieve categories at the moment",
// 		}.RenderError(w)
// 		return
// 	}

// userID, ok := r.Context().Value(models.UserIDKey).(int)
// if !ok {
// 	models.Error{
// 		User:       user,
// 		StatusCode: http.StatusInternalServerError,
// 		Message:    "Internal Server Error",
// 		SubMessage: "Unable to retrieve user information.",
// 	}.RenderError(w)
// 	return
// }

// if userID != 0 {
// 	user, err = webForum.Users.FindUserByID(userID)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			user = &models.User{
// 				UserType: "guest",
// 			}
// 		} else {
// 			models.Error{
// 				User:       user,
// 				StatusCode: http.StatusInternalServerError,
// 				Message:    "Internal Server Error",
// 				SubMessage: err.Error(),
// 			}.RenderError(w)
// 			return
// 		}
// 	}
// } else {
// 	user = &models.User{
// 		UserType: "guest",
// 	}
// }
// userType, ok := r.Context().Value(userTypeKey).(string)
// if !ok {
// 	models.Error{
// 		User:       user,
// 		StatusCode: http.StatusInternalServerError,
// 		Message:    "Internal Server Error",
// 		SubMessage: "Unable to retrieve user type.",
// 	}.RenderError(w)
// 	return
// }
// user.UserType = userType

// 	CreatePostData := struct {
// 		User       *models.User
// 		Categories []models.Category
// 	}{
// 		User:       user,
// 		Categories: categories,
// 	}
// 	models.RenderPage(w, "post-creation.html", CreatePostData)
// }

// func (webForum *WebApp) Creation(w http.ResponseWriter, r *http.Request) {
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

// 	var postData PostUpdate
// 	err = json.NewDecoder(r.Body).Decode(&postData)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 		return
// 	}

// 	if postData.Title == "" || len(postData.Title) > 70 || postData.Content == "" || len(postData.Content) > 5000 || len(postData.Categories) == 0 || len(postData.Categories) > 3 {
// 		http.Redirect(w, r, "/createPost", http.StatusSeeOther)
// 		return
// 	}

// 	ids, err := webForum.Post.GetIdsCategories(postData.Categories)
// 	if err != nil {
// 		models.Error{
// 			User:       &models.User{},
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "500 Internal Server Error",
// 			SubMessage: "Oops! " + err.Error(),
// 		}.RenderError(w)
// 		return
// 	}

// 	idPost, err := webForum.Post.CreatePost(userID, postData.Title, postData.Content)
// 	if err != nil {
// 		models.Error{
// 			User:       &models.User{},
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "500 Internal Server Error",
// 			SubMessage: "Oops! " + err.Error(),
// 		}.RenderError(w)
// 		return
// 	}

// 	err = webForum.Post.AddCategoriesPost(idPost, ids)
// 	if err != nil {
// 		models.Error{
// 			User:       &models.User{},
// 			StatusCode: http.StatusInternalServerError,
// 			Message:    "500 Internal Server Error",
// 			SubMessage: "Oops! " + err.Error(),
// 		}.RenderError(w)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusSeeOther)
// }

// func (webForum *WebApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

// 	post, err := webForum.Post.UpdatePost(userID, id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	Categorys, err := webForum.Post.GetCategories()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	data := struct {
// 		User       *models.User
// 		Post_info  models.Post
// 		Categories []models.Category
// 	}{
// 		User:       user,
// 		Post_info:  post,
// 		Categories: Categorys,
// 	}

// 	models.RenderPage(w, "post-update.html", data)
// }

// type PostUpdate struct {
// 	Id         string   `json:"id`
// 	Title      string   `json:"title"`
// 	Content    string   `json:"content"`
// 	Categories []string `json:"categories"`
// }

// func (webForum *WebApp) Updating(w http.ResponseWriter, r *http.Request) {
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

// // madara
type FetchComments struct {
	PostId string `json:"ID"`
}

func (webForum *WebApp) GetComments(w http.ResponseWriter, r *http.Request) {
	var commentData FetchComments

	if decodeErr := decodeJsonData(r, &commentData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// fmt.Println("postId:", commentData.PostId)
	PostID, err := strconv.Atoi(commentData.PostId)
	if err != nil {
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}
	comments, err := webForum.Comments.Comments(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("comments:", comments)
	if err := encodeJsonData(w, http.StatusOK, comments); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

type CreateComment struct {
	PostID  string `json:"postid"`
	// UserID  string `json:"userid"`
	Content string `json:"content"`
}

func (webForum *WebApp) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}
	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newCmntData CreateComment
	if decodeErr := decodeJsonData(r, &newCmntData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(newCmntData.PostID)
	if err != nil {
		http.Error(w, "Invalid Post ID ", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newCmntData.Content) == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}
	err = webForum.Comments.CreateComment(postID, user.ID, newCmntData.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, "comments created succesfully"); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// // tools-------------------------------------

func decodeJsonData(r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func encodeJsonData(w http.ResponseWriter, statusCode int, obj interface{}) error {
	// fmt.Println("OBJ:\x1b[1;31m", obj, "\x1b[1;39m")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}

	return nil
}
