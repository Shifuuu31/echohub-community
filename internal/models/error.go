package models

import (
	"html/template"
	"net/http"
)

var Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))

type Error struct {
	User       *User `json:"user"`
	StatusCode int  `json:"statusCode"`
	Message    string `json:"message"`
	SubMessage string `json:"subMessage"`
	Type string `json:"type"`
}

func (err Error) RenderError(w http.ResponseWriter) {
	if err := Template.ExecuteTemplate(w, "error.html", err); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderPage(w http.ResponseWriter, name string, obj interface{}) {
	if err := Template.ExecuteTemplate(w, name, obj); err != nil {
		return
	}
}
