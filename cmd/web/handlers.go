package web

import (
	"context"
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
	// var err error
	

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
		user, err := webForum.Users.FindUserByID(userID)
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
		print("HAHAHAHAHAHAH\n")
		err := models.Error{
			User:       user, // Ensure this is not nil
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! The page you are looking for does not exist",
		}
		err.RenderError(w)
		return
	}

	models.RenderPage(w, "home.html", user)
}

func (webForum *WebApp) LoginPage(w http.ResponseWriter, r *http.Request) {
	models.RenderPage(w, "login.html", nil)
}

func (webForum *WebApp) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "A login error occurred",
		}.RenderError(w)
		return
	}

	userID, err := webForum.Users.ValidateUserCredentials(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		models.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
			SubMessage: "Invalid username or password",
		}.RenderError(w)
		return
	}

	newSession, err := webForum.Sessions.GenerateNewSession(userID, r.FormValue("remember"))
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Cannot generate new session",
		}.RenderError(w)
		return
	}

	newCookie, err := webForum.Sessions.InsertOrUpdateSession(newSession)
	if err != nil {
		models.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			SubMessage: "Cannot insert new session",
		}.RenderError(w)
		return
	}
	http.SetCookie(w, &newCookie)

	http.Redirect(w, r, "/", http.StatusFound)
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

	http.Redirect(w, r, "/", http.StatusFound)
}

func (webForum *WebApp) RegisterPage(w http.ResponseWriter, r *http.Request) {
	models.RenderPage(w, "register.html", nil)
}

func (webForum *WebApp) UserRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "a registration error occured"}.RenderError(w)
		return
	}

	newUser, err := webForum.Users.ValidateNewUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), r.FormValue("rPassword"))
	if err != nil {
		models.Error{StatusCode: http.StatusBadRequest, Message: "Bad Request", SubMessage: "Invalid input data"}.RenderError(w) // ,have to handle bcrypt error as internal server error
		return
	}

	if err := webForum.Users.InsertUser(newUser); err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot insert new user"}.RenderError(w)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
