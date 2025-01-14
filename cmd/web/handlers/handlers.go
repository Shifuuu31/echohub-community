package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"forum/internal/models"
)

var (
	Template       = template.Must(template.ParseGlob("./cmd/web/templates/*.html"))
	categoriesForm = []string{}
)

func (webForum *WebApp) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	categories, err := webForum.Post.GetCategorys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	Categories_Posts, err := webForum.Post.GetCategoriesNames(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type data struct {
		ID               int
		Categories       []models.Categorie
		Username         string
		Title            string
		Post_content     string
		Categories_Posts []string
		Creation_date    time.Time
	}
	var datas []data
	for i := range posts {
		datas = append(datas, data{
			ID:               posts[i].ID,
			Categories:       categories,
			Username:         users[i],
			Title:            posts[i].Title,
			Post_content:     posts[i].Post_content,
			Categories_Posts: Categories_Posts[i],
			Creation_date:    posts[i].Creation_date,
		})
	}

	if err := Template.ExecuteTemplate(w, "index.html", datas); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) CreatePostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	categories, err := webForum.Post.GetCategorys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-create.html", categories); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
		return
	}
}

func (webForm *WebApp) NewPostCreationHandler(w http.ResponseWriter, r *http.Request) {
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

func (webForum *WebApp) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/delete" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	// fmt.Printf("%V",id)
	err = webForum.Post.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (WebForum *WebApp) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/Update" {
		w.WriteHeader(http.StatusNotFound)
		if err := Template.ExecuteTemplate(w, "404.html", nil); err != nil {
			http.Error(w, "Error loading 404 Page", http.StatusInternalServerError)
			return
		}
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("ID"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}

	err = WebForum.Post.UpdatPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := Template.ExecuteTemplate(w, "post-create.html", nil); err != nil {
		http.Error(w, "Error loading HomePage", http.StatusInternalServerError)
	}
}
