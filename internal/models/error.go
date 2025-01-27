package models

import (
	"html/template"
	"net/http"
)

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

type Error struct {
	User       *User
	StatusCode int
	Message    string
	SubMessage string
}

func (err Error) RenderError(w http.ResponseWriter) {
	if err := Template.ExecuteTemplate(w, "error.html", err); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderPage(w http.ResponseWriter, name string, obj interface{}) {
	if err := Template.ExecuteTemplate(w, name, obj); err != nil {
		Error{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error", SubMessage: "Error loading page"}.RenderError(w)
		return
	}
}
