package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	fmt.Println(user)
	homeData := struct {
		User *models.User
	}{
		User: user,
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
