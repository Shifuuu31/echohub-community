package handlers

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"

	"echohub-community/internal/models"
)

// HomePage renders the home page
// @Summary      Show home page
// @Description  Render the home page with user info and categories
// @Tags         Home
// @Produce      html
// @Success      200  {string}  string  "Home page HTML"
// @Router       / [get]
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
	homeData.User.UserName = html.EscapeString(homeData.User.UserName)

	models.RenderPage(w, "home.html", homeData)
}

// MaxID returns the maximum post ID
// @Summary      Get max post ID
// @Description  Get the highest post ID for pagination/client-side tracking
// @Tags         Posts
// @Produce      json
// @Success      200  {integer}  int
// @Failure      500  {object}   models.Error
// @Router       /maxId [post]
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

// GetPosts retrieves posts with pagination and category filter
// @Summary      Get posts
// @Description  Retrieve a list of posts based on start ID and category
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        postsData  body   object  true  "Post filter"  example({"start":0,"category":"all"})
// @Success      200  {array}   models.Post
// @Failure      400  {string}  string  "Invalid JSON format"
// @Failure      500  {object}  models.Error
// @Router       /posts [post]
func (webForum *WebApp) GetPosts(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	postsData := struct {
		StartId  int    `json:"start"`
		Category string `json:"category"`
	}{}

	if decodeErr := decodeJsonData(r, &postsData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	posts, postErr := webForum.Posts.GetPosts(user.ID, postsData.StartId, postsData.Category)
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

// GetComments retrieves comments for a specific post
// @Summary      Get comments
// @Description  Retrieve all comments for a post by its ID
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        commentData  body   FetchComments  true  "Comment fetch data"
// @Success      200  {array}   models.Comment
// @Failure      400  {string}  string  "Invalid JSON format or PostID"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /comments [post]
func (webForum *WebApp) GetComments(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		userErr.RenderError(w)
		return
	}

	var commentData FetchComments

	if decodeErr := decodeJsonData(r, &commentData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	PostID, err := strconv.Atoi(commentData.PostId)
	if err != nil {
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}
	comments, cmtErr := webForum.Comments.GetPostComments(PostID, user.ID)
	if cmtErr.Type == "server" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := encodeJsonData(w, http.StatusOK, comments); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

type CreateComment struct {
	PostID string `json:"postid"`
	// UserID  string `json:"userid"`
	Content string `json:"content"`
}

// CreateComment adds a new comment to a post
// @Summary      Create comment
// @Description  Post a new comment to an existing post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        newCmntData  body   CreateComment  true  "New comment data"
// @Success      200  {string}  string  "comments created succesfully"
// @Failure      400  {string}  string  "Invalid JSON or content"
// @Failure      403  {string}  string  "Forbidden"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /createComment [post]
func (webForum *WebApp) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}
	if user.UserType != "authenticated" {
		if err := encodeJsonData(w, http.StatusForbidden, "Forbidden"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	var newCmntData CreateComment
	if decodeErr := decodeJsonData(r, &newCmntData); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	postID, err := strconv.Atoi(newCmntData.PostID)
	if err != nil {
		if err := encodeJsonData(w, http.StatusBadRequest, "Invalid Post ID"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	if strings.TrimSpace(newCmntData.Content) == "" {
		if err := encodeJsonData(w, http.StatusBadRequest, "Comment cannot be empty"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}
	err = webForum.Comments.CreateComment(postID, user.ID, newCmntData.Content)
	if err != nil {
		if err := encodeJsonData(w, http.StatusInternalServerError, "Failed to create comment"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}

	if err := encodeJsonData(w, http.StatusOK, "comments created succesfully"); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// LikeDislikeHandler (when we click on like/dislike button)
// LikeDislikeHandler toggles a like or dislike on a post or comment
// @Summary      Like/Dislike
// @Description  Add or remove a like/dislike from a post or comment
// @Tags         Interactions
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        request  body   models.Interaction  true  "Interaction data"
// @Success      200  {object}  models.Response
// @Failure      400  {string}  string  "Invalid entity type or JSON"
// @Failure      403  {string}  string  "Forbidden"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /like-dislike [post]
func (webForum *WebApp) LikeDislikeHandler(w http.ResponseWriter, r *http.Request) {
	user, userErr := webForum.Users.RetrieveUser(r)
	if userErr.Type == "server" {
		http.Error(w, userErr.Message, userErr.StatusCode)
		return
	}
	if user.UserType != "authenticated" {
		if err := encodeJsonData(w, http.StatusForbidden, "Forbidden: please try to login"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}
	var request models.Interaction
	if decodeErr := decodeJsonData(r, &request); decodeErr != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if request.EntityType != "post" && request.EntityType != "comment" {
		if err := encodeJsonData(w, http.StatusBadRequest, "Invalid entity type"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		return
	}
	if err := webForum.LikesDislikes.LikeDislike(request.EntityID, request.EntityType, user.ID, request.Liked); err != nil {
		if err := encodeJsonData(w, http.StatusBadRequest, "Invalid entity type"); err != nil {
			http.Error(w, "failed to encode object.", http.StatusInternalServerError)
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	reaction, reactErr := models.GetReaction(webForum.LikesDislikes.DB, request.EntityID, request.EntityType, user.ID)
	if reactErr.Type == "server" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var response models.Response
	response.Messages = append(response.Messages, reaction)
	likeCount, dislikeCount, ldErr := models.GetLikesDislikesCount(webForum.LikesDislikes.DB, request.EntityID, request.EntityType)
	if ldErr.Type == "server" {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response.Extra = append(response.Extra, strconv.Itoa(likeCount))
	response.Extra = append(response.Extra, strconv.Itoa(dislikeCount))
	if err := encodeJsonData(w, http.StatusOK, response); err != nil {
		http.Error(w, "failed to encode object.", http.StatusInternalServerError)
	}
}

// // tools-------------------------------------

func decodeJsonData(r *http.Request, obj interface{}) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

func encodeJsonData(w http.ResponseWriter, statusCode int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return err
	}

	return nil
}
