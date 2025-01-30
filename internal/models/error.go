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
	Type string
}

func (err Error) RenderError(w http.ResponseWriter) {
	if err := Template.ExecuteTemplate(w, "error.html", err); err != nil {
		print("hjgcysgycs")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderPage(w http.ResponseWriter, name string, obj interface{}) {
	if err := Template.ExecuteTemplate(w, name, obj); err != nil {
		return
	}
}
