package handlers

import (
	"context"
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
	models.RenderPage(w, "profileSettings.html", user)
}

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
		fmt.Println("KO")
	} else {
		fmt.Println("OK")
		response.Messages = append(response.Messages, "Profile Updated successfully")
		if err := encodeJsonData(w, http.StatusOK, response); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
	}
}
