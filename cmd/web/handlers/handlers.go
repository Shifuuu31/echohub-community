package handlers

import (
	"net/http"
	"text/template"
)

func (webForum *WebApp) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	HomePage, err := template.ParseFiles("cmd/web/templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ErrorPage, err := template.ParseFiles("cmd/web/templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := webForum.Post.GetPosts()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		ErrorPage.Execute(w, nil)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed ", http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
	err = HomePage.Execute(w, posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
