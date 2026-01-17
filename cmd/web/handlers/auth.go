package handlers

import (
	"context"
	"fmt"
	"html"
	"net/http"

	"echohub-community/internal/models"
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

// LoginPage renders the login page
// @Summary      Show login page
// @Description  Render the login page for the user
// @Tags         Auth
// @Produce      html
// @Success      200  {string}  string  "Login page HTML"
// @Router       /login [get]
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

// ConfirmLogin handles user authentication
// @Summary      User login
// @Description  Authenticate user with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body   object  true  "Login credentials"  example({"username":"user1","password":"password123","rememberMe":true})
// @Success      200  {object}  models.Response
// @Failure      400  {string}  string  "Invalid JSON format"
// @Failure      401  {object}  models.Response  "Unauthorized"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /confirmLogin [post]
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

// UserLogout logs out the current user
// @Summary      User logout
// @Description  Expire the user session and logout
// @Tags         Auth
// @Success      302  {string}  string  "Redirect to /login"
// @Router       /logout [get]
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

// RegisterPage renders the registration page
// @Summary      Show registration page
// @Description  Render the registration page for the user
// @Tags         Auth
// @Produce      html
// @Success      200  {string}  string  "Registration page HTML"
// @Router       /register [get]
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

// UserRegister handles new user registration
// @Summary      User registration
// @Description  Register a new user with email, username, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        newUser  body   models.NewUserInfo  true  "User registration data"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response  "Validation errors"
// @Failure      500  {object}  models.Error  "Internal server error"
// @Router       /confirmRegister [post]
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

// ProfileSettings renders the profile settings page
// @Summary      Show profile settings
// @Description  Render the profile settings page for authenticated users
// @Tags         User
// @Produce      html
// @Security     CookieAuth
// @Success      200  {string}  string  "Profile settings page HTML"
// @Failure      401  {object}  models.Error  "Unauthorized"
// @Router       /profileSettings [get]
func (webForum *WebApp) ProfileSettings(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if user.UserType != "authenticated" {
		userErr = models.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			SubMessage: "Try to login",
		}
		userErr.RenderError(w)
		return
	}
	user.UserName = html.EscapeString(user.UserName)

	models.RenderPage(w, "profileSettings.html", user)
}

// UpdateProfile updates the user's profile information
// @Summary      Update profile
// @Description  Update user's nickname, email, or other details
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        toUpdate  body   models.NewUserInfo  true  "Updated profile data"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response  "Validation errors"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {object}  models.Error  "Internal server error"
// @Router       /updateProfile [post]
func (webForum *WebApp) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}

	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var toUpdate models.NewUserInfo

	if decodeErr := decodeJsonData(r, &toUpdate); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Println(toUpdate)
	response, err := webForum.Users.UpdateUser(toUpdate, user.ID)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot update new user"}.RenderError(w)
		return
	}

	if len(response.Messages) != 0 {
		if err := encodeJsonData(w, http.StatusBadRequest, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	} else {
		response.Messages = append(response.Messages, "Profile Updated successfully")
		if err := encodeJsonData(w, http.StatusOK, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	}
}
