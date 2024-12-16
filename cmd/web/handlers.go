package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
)

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return

	}
	if err := Template.ExecuteTemplate(w, "home-guestMode.html", nil); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForum *WebApp) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return

	}
	if err := Template.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Error loading LoginPage", http.StatusInternalServerError)
		return
	}
}

func (webForum *WebApp) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "a registration error occured", http.StatusInternalServerError)
		return
	}

	userID, err := webForum.Users.ValidateUserCreadentials(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(r.Context(), "userID", userID)

	webForum.SetCookie(w, r.WithContext(ctx))
}

func (webForum *WebApp) SetCookie(w http.ResponseWriter, r *http.Request) {
	value := r.Context().Value("userID")

	userID, ok := value.(int)
	if !ok {
		http.Error(w, "No user found in context", http.StatusInternalServerError)
		return
	}

	fmt.Printf("sCookie: UserID: %d\n", userID)

	newSession, err := webForum.Sessions.GenerateNewSession(userID)
	if err != nil {
		http.Error(w, "cannot generate new session", http.StatusInternalServerError)
		return
	}
	fmt.Println(newSession)

	newCookie, err := webForum.Sessions.InsertSession(newSession)
	if err != nil {
		http.Error(w, "cannot insert new session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &newCookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (webForum *WebApp) GetCookie(w http.ResponseWriter, r *http.Request) {
}

func (webForum *WebApp) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return

	}
	if err := Template.ExecuteTemplate(w, "register.html", nil); err != nil {
		http.Error(w, "Error loading LoginPage", http.StatusInternalServerError)
		return
	}
}

func (webForum *WebApp) UserRegister(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "a registration error occured", http.StatusInternalServerError)
		return
	}

	newUser, err := webForum.Users.ValidateNewUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"), r.FormValue("rPassword"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := webForum.Users.InsertUser(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
