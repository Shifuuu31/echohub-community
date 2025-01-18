package web

import (
	"context"
	"fmt"
	"net/http"
	"forum/internal/models"
)

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}
	models.RenderPage(w, "home-guestMode.html", nil)
}

func (webForum *WebApp) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}
	models.RenderPage(w, "login.html", nil)
}

func (webForum *WebApp) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "A login error occured"}.RenderError(w)
		return
	}

	userID, err := webForum.Users.ValidateUserCreadentials(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		models.Error{StatusCode: http.StatusUnauthorized, Message: "Unauthorized", SubMessage: "Invalid username or password"}.RenderError(w)
		return
	}
	ctx := context.WithValue(r.Context(), "UserID", userID)

	webForum.SetCookie(w, r.WithContext(ctx))
}

func (webForum *WebApp) SetCookie(w http.ResponseWriter, r *http.Request) {
	value := r.Context().Value("UserID")

	userID, ok := value.(int)
	if !ok {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "No user found in context"}.RenderError(w)
		// http.Error(w, "No user found in context", http.StatusInternalServerError)
		return
	}

	// fmt.Printf("sCookie: UserID: %d\n", userID)

	newSession, err := webForum.Sessions.GenerateNewSession(userID)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot generate new session"}.RenderError(w)
		// http.Error(w, "cannot generate new session", http.StatusInternalServerError)
		return
	}
	fmt.Println(newSession)

	newCookie, err := webForum.Sessions.InsertSession(newSession)
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot insert new session"}.RenderError(w)

		// http.Error(w, "cannot insert new session", http.StatusInternalServerError)
		return
	}
	// fmt.Println(newCookie)
	http.SetCookie(w, &newCookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (webForum *WebApp) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return

	}
	models.RenderPage(w, "register.html", nil)
}

func (webForum *WebApp) UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		models.Error{StatusCode: http.StatusNotFound, Message: "404 Page Not Found", SubMessage: "Oops! the page you looking for does not exist"}.RenderError(w)
		return
	}

	err := r.ParseForm()
	if err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "a registration error occured"}.RenderError(w)
		return
	}

	newUser, err := webForum.Users.ValidateNewUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), r.FormValue("rPassword"))
	if err != nil {
		models.Error{StatusCode: http.StatusBadRequest, Message: "Bad Request", SubMessage: "Invalid input data"}.RenderError(w) // ,have to handle bcrypt error as internal server error
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := webForum.Users.InsertUser(newUser); err != nil {
		models.Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Cannot insert new user"}.RenderError(w)

		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
