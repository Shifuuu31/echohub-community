package handlers

import (
	"net/http"
	"text/template"
)

func HandleHomePage(w http.ResponseWriter, r *http.Request) {

	is_login := false
	var HomePage *template.Template
	var err error

	if is_login {
		HomePage, err = template.ParseFiles("cmd/web/templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		HomePage, err = template.ParseFiles("cmd/web/templates/home2.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	ErrorPage, err := template.ParseFiles("cmd/web/templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		ErrorPage.Execute(w, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = HomePage.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
