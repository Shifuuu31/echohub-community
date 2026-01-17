package models

import (
	"html/template"
	"net/http"
)

var Template *template.Template

func LoadTemplates() {
	if Template == nil {
		Template = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))
	}
}

type Error struct {
	User       *User  `json:"User"`
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message"`
	SubMessage string `json:"SubMessage"`
	Type       string `json:"Type"`
}

func (err Error) RenderError(w http.ResponseWriter) {
	if Template == nil {
		http.Error(w, err.Message, err.StatusCode)
		return
	}
	if err := Template.ExecuteTemplate(w, "error.html", err); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderPage(w http.ResponseWriter, name string, obj interface{}) {
	if Template == nil {
		http.Error(w, "Internal Server Error: Templates not loaded", http.StatusInternalServerError)
		return
	}
	if err := Template.ExecuteTemplate(w, name, obj); err != nil {
		return
	}
}
