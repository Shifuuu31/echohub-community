package handlers

import (
	"html/template"
	"net/http"
	"time"

	"forum/internal/models"
)

var (
	Template       = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))
	categoriesForm = []string{}
)

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	posts, err := webForum.Post.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := webForum.Post.GetUsersNames(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := webForum.Post.GetCategoriesNames(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type data struct {
		Username      string
		Title         string
		Post_content  string
		Categories    []string
		Creation_date time.Time
	}
	datas := []data{}
	for i := 0; i < len(posts); i++ {
		data := data{}
		data.Username = users[i]
		data.Title = posts[i].Title
		data.Post_content = posts[i].Post_content
		data.Categories = categories[i]
		data.Creation_date = posts[i].Creation_date
		datas = append(datas, data)
	}

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	if err := Template.ExecuteTemplate(w, "index.html", datas); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) CreatePostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	categories, err := webForm.Post.GetCategorys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-create.html", categories); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) NewPostCreation(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	New := models.Post{
		Title:        r.FormValue("title"),
		Post_content: r.FormValue("content"),
	}
	categoriesForm = r.Form["categories[]"]

	ids, err := webForm.Post.GetIdsCategories(categoriesForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = webForm.Post.CreatePost(New.Title, New.Post_content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	idPost, err := webForm.Post.GetLastPoID()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = webForm.Post.AddCategoriePost(idPost, ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
