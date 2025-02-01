package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
		http.Error(w, "failed to decode object.", http.StatusInternalServerError)
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
		http.Error(w, "failed to decode object.", http.StatusInternalServerError)
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
		if err := encodeJsonData(w, maxIdErr.StatusCode, maxIdErr); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	if err := encodeJsonData(w, http.StatusOK, maxID); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	postsData := struct {
		StartId  int    `json:"postID"`
		Category string `json:"category"`
	}{}

	if decodeErr := decodeJsonData(r, &postsData); decodeErr != nil {
		http.Error(w, "failed to decode object.", http.StatusInternalServerError)
		return
	}

	posts, postErr := webForum.Post.GetPosts(postsData.StartId, postsData.Category)
	fmt.Println(len(posts))
	if len(posts) == 0 {
		fmt.Println(posts)
	}
	if postErr.Type != "" {
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
		fmt.Println("here")
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
// type FetchComments struct {
// 	Id string `json:ID`
// }

// func (webForum *WebApp) GetComments(w http.ResponseWriter, r *http.Request) {
// 	// w.Header().Set("Access-Control-Allow-Credentials", "false")
// 	log.Println(r.Cookie("userSession"))
// 	var commentData FetchComments
// 	err := json.NewDecoder(r.Body).Decode(&commentData)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 		return
// 	}

// 	PostId, err := strconv.Atoi(commentData.Id)
// 	fmt.Println(PostId)
// 	if err != nil {
// 		return
// 	}
// 	commentModel := &models.CommentModel{DB: webForum.Post.DB}
// 	comments, err := commentModel.GetComments(PostId)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	encodeJsonData(w, comments)
// }

// type CreateComment struct {
// 	PostID         string `json:"postid"`
// 	UserID         int    `json:"userid"`
// 	CommentContent string `json:"content"`
// }

// func (webForum *WebApp) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
// 	var commentData CreateComment
// 	err := json.NewDecoder(r.Body).Decode(&commentData)
// 	fmt.Println(commentData)
// 	if err != nil {
// 		fmt.Fprint(w, "HAHHAHAHHAHA")
// 		// http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		fmt.Fprint(w, "11111")
// 		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// 	postID, err := strconv.Atoi(commentData.PostID)
// 	userID := commentData.UserID
// 	if err != nil {
// 		fmt.Fprint(w, "222222")
// 		// http.Error(w, "Invalid post ID or user ID", http.StatusBadRequest)
// 		return
// 	}

// 	content := commentData.CommentContent
// 	if content == "" {
// 		fmt.Fprint(w, "333333")
// 		// http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
// 		return
// 	}
// 	commentModel := &models.CommentModel{DB: webForum.Post.DB}
// 	err = commentModel.CreateComment(postID, userID, content)
// 	if err != nil {
// 		fmt.Fprint(w, "444444")
// 		// http.Error(w, "Failed to create comment", http.StatusInternalServerError)
// 		return
// 	}
// 	// fmt.Fprintln(w, "comment created")
// 	encodeJsonData(w, "{}")
// }

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
