package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
)

func (webForum *WebApp) HomePage(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	if r.URL.Path != "/" {
		models.Error{
			User:       user,
			StatusCode: http.StatusNotFound,
			Message:    "404 Page Not Found",
			SubMessage: "Oops! The page you are looking for does not exist",
		}.RenderError(w)
		return
	}

	categories, catsErr := webForum.Posts.GetCategories()
	if catsErr.Type == "server" {
		catsErr.RenderError(w)
		return
	}

	homeData := struct {
		User       *models.User
		Categories []models.Category
	}{
		User:       user,
		Categories: categories,
	}

	models.RenderPage(w, "home.html", homeData)
}

func (webForum *WebApp) MaxID(w http.ResponseWriter, r *http.Request) {
	maxID, maxIdErr := webForum.Posts.GetMaxId()
	if maxIdErr.Type == "server" {
		maxIdErr.RenderError(w)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, maxID); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

func (webForum *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	postsData := struct {
		StartId  int    `json:"start"`
		Category string `json:"category"`
	}{}

	if decodeErr := decodeJsonData(r, &postsData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	posts, postErr := webForum.Posts.GetPosts(postsData.StartId, postsData.Category)
	if postErr.Type != "" {
		if postErr.Type == "server" {
			postErr.RenderError(w)
			return
		}

		if err := encodeJsonData(w, postErr.StatusCode, postErr); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	if len(posts) == 0 {
		postErr = models.Error{
			StatusCode: http.StatusContinue,
			Message:    "No posts available",
			Type:       "client",
		}

		if err := encodeJsonData(w, postErr.StatusCode, postErr); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return

	}

	if err := encodeJsonData(w, http.StatusOK, posts); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}




// // madara
type FetchComments struct {
	PostId string `json:"ID"`
}

func (webForum *WebApp) GetComments(w http.ResponseWriter, r *http.Request) {
	var commentData FetchComments

	if decodeErr := decodeJsonData(r, &commentData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// fmt.Println("postId:", commentData.PostId)
	PostID, err := strconv.Atoi(commentData.PostId)
	if err != nil {
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}
	comments, err := webForum.Comments.Comments(PostID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("comments:", comments)
	if err := encodeJsonData(w, http.StatusOK, comments); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

type CreateComment struct {
	PostID string `json:"postid"`
	// UserID  string `json:"userid"`
	Content string `json:"content"`
}

func (webForum *WebApp) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}
	if user.UserType != "authenticated" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newCmntData CreateComment
	if decodeErr := decodeJsonData(r, &newCmntData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(newCmntData.PostID)
	if err != nil {
		http.Error(w, "Invalid Post ID ", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newCmntData.Content) == "" {
		http.Error(w, "Comment cannot be empty", http.StatusBadRequest)
		return
	}
	err = webForum.Comments.CreateComment(postID, user.ID, newCmntData.Content)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	if err := encodeJsonData(w, http.StatusOK, "comments created succesfully"); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// // tools-------------------------------------

func decodeJsonData(r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func encodeJsonData(w http.ResponseWriter, statusCode int, obj interface{}) error {
	// fmt.Println("OBJ:\x1b[1;31m", obj, "\x1b[1;39m")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}

	return nil
}
