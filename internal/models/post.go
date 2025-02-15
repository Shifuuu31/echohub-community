package models

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID            int
	UserName      string
	ProfileImg    string
	CreatedAt     time.Time
	Title         string
	Content       string
	Categories    []string
	LikeCount     int
	DislikeCount  int
	CommentsCount int
	Reaction      string
}

type Category struct {
	ID               int
	CategoryName     string
	CategoryIconPath string
}

type PostModel struct {
	DB *sql.DB
}

// get categories from DB
func (PostModel *PostModel) GetCategories() (Categories []Category, catsErr Error) {
	selectStmt := `	SELECT
					    id,
					    category_name,
					    category_icon_path
					FROM
					    Categories;`

	catsErr = Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		SubMessage: "Unable to get categories.",
		Type:       "server",
	}

	rowsDB, err := PostModel.DB.Query(selectStmt)
	if err != nil {
		return nil, catsErr
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		category := Category{}
		if err := rowsDB.Scan(&category.ID, &category.CategoryName, &category.CategoryIconPath); err != nil {
			return nil, catsErr
		}
		Categories = append(Categories, category)
	}

	if err = rowsDB.Err(); err != nil {
		return nil, catsErr
	}

	return Categories, Error{}
}

// Get maxID of posts
func (postModel *PostModel) GetMaxId() (maxID int, maxIdError Error) {
	if err := postModel.DB.QueryRow("SELECT p.id FROM PostTable p ORDER BY p.id DESC LIMIT 1").Scan(&maxID); err != nil {
		if err == sql.ErrNoRows {
			return maxID, maxIdError
		}
		return maxID, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	return maxID, maxIdError
}

// get posts from DB with cateogry
func (postModel *PostModel) GetPosts(userID int, startId int, category string) (posts []Post, postsErr Error) {
	var (
		query      string
		args       []interface{}
		categoryID int
	)

	switch category {
	case "All":
		query = `SELECT
					PostTable.id,
					UserTable.username,
					UserTable.profile_img,
					PostTable.title,
					PostTable.content,
					PostTable.creation_date
				FROM
					PostTable
					JOIN UserTable ON UserTable.id = PostTable.user_id
				WHERE
					PostTable.id <= ?
				ORDER BY
					PostTable.id DESC
				LIMIT
					10;`
		args = append(args, startId)

	case "MyPosts":
		query = `SELECT
					PostTable.id,
					UserTable.username,
					UserTable.profile_img,
					PostTable.title,
					PostTable.content,
					PostTable.creation_date
				FROM
					PostTable
					JOIN UserTable ON UserTable.id = PostTable.user_id
				WHERE
					UserTable.id = ? AND PostTable.id <= ?
				ORDER BY
					PostTable.id DESC
				LIMIT
					10;`
		args = append(args, userID, startId)

	case "LikedPosts":
		// Placeholder for "LikedPosts" case
		query = `SELECT
					PostTable.id,
					UserTable.username,
					UserTable.profile_img,
					PostTable.title,
					PostTable.content,
					PostTable.creation_date
				FROM
					PostTable
					JOIN UserTable ON UserTable.id = PostTable.user_id
					JOIN Likes_Dislikes ON Likes_Dislikes.entity_id = PostTable.id
					AND Likes_Dislikes.entity_type = "post"
				WHERE
					PostTable.id <= ?
					AND Likes_Dislikes.user_id = ?
				ORDER BY
					PostTable.id DESC
				LIMIT
					10;`
		args = append(args, startId, userID)
	default:
		var err error
		if categoryID, err = strconv.Atoi(category); err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid category",
				Type:       "client",
			}
		}

		query = `SELECT
					PostTable.id,
					UserTable.username,
					UserTable.profile_img,
					PostTable.title,
					PostTable.content,
					PostTable.creation_date
				FROM
					PostTable
					JOIN Categories_Posts ON PostTable.id = Categories_Posts.post_id
					JOIN UserTable ON UserTable.id = PostTable.user_id
				WHERE
					Categories_Posts.category_id = ?
					AND PostTable.id <= ?
				ORDER BY
					PostTable.id DESC
				LIMIT
					10;`
		args = append(args, categoryID, startId)
	}

	rows, err := postModel.DB.Query(query, args...)
	if err != nil {
		return []Post{}, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	defer rows.Close()

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.ID, &post.UserName, &post.ProfileImg, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}

		if post.Categories, err = postModel.GetPostCategories(post.ID); err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}
		post.CreatedAt = post.CreatedAt.UTC()
		post.CommentsCount, err = postModel.GetCommentCount(post.ID)
		if err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}

		var LDErr Error
		if post.LikeCount, post.DislikeCount, LDErr = GetLikesDislikesCount(postModel.DB, post.ID, "post"); LDErr.Message != "" {
			return []Post{}, LDErr
		}
		if userID > 0 {
			Reaction, ReactionErr := GetReaction(postModel.DB, post.ID, "post", userID)
			if ReactionErr.Message != "" {
				return []Post{}, LDErr
			}
			post.Reaction = Reaction
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return []Post{}, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}

	return posts, postsErr
}

// get cateogries of post
func (post *PostModel) GetPostCategories(postId int) (postCategories []string, err error) {
	query := `SELECT
			      Categories.category_name
			  FROM
			      Categories_Posts
			      JOIN Categories ON Categories_Posts.category_id = Categories.id
			  WHERE
			      Categories_Posts.post_id = ?;`

	rows, err := post.DB.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryName string
		if err := rows.Scan(&categoryName); err != nil {
			return nil, err
		}
		postCategories = append(postCategories, categoryName)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return postCategories, nil
}

type PostData struct {
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"selectedCategories"`
}

func (post *PostModel) GetCommentCount(postID int) (cmntCount int, err error) {
	countStmt := `SELECT COUNT(*) FROM CommentTable WHERE post_id = ?`

	err = post.DB.QueryRow(countStmt, postID).Scan(&cmntCount)
	return cmntCount, err
}

func CheckNewPost(newPost PostData) (response Response) {
	if len(newPost.Categories) == 0 {
		response.Messages = append(response.Messages, "Select at least one category")
	}
	if len(newPost.Categories) > 3 {
		response.Messages = append(response.Messages, "You can only select up to 3 categories")
	}
	if strings.TrimSpace(newPost.Title) == "" {
		response.Messages = append(response.Messages, "Title cannot be empty")
	}

	if len(newPost.Title) > 70 {
		response.Messages = append(response.Messages, "Title length up to 70 character")
	}

	if len(newPost.Content) > 5000 {
		response.Messages = append(response.Messages, "Content length up to 5000 character")
	}
	return response
}

// create post (insert in DB)
func (PostModel *PostModel) CreatePost(userId int, title, content string) (int, error) {
	var id int
	query := `	INSERT INTO
				    PostTable (title, user_id, content)
				VALUES
				    (?, ?, ?) RETURNING id;`
	err := PostModel.DB.QueryRow(query, title, userId, content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// get post by id to update
func (PostModel *PostModel) GetPost(user_id, idPost int) (post Post, postErr Error) {
	err := PostModel.DB.QueryRow("SELECT id, title, content FROM PostTable WHERE id = ? AND user_id = ?", idPost, user_id).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, Error{
				StatusCode: http.StatusForbidden,
				Message:    "Forbidden",
				Type:       "client",
			}
		}
		return Post{}, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}

	if post.Categories, err = PostModel.GetPostCategories(post.ID); err != nil {
		return Post{}, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}

	return post, Error{}
}

// update post
func (PostModel *PostModel) EditPost(idPost int, title string, content string, Categories []string) (err error) {
	_, err = PostModel.DB.Exec("UPDATE PostTable SET title = ?, content = ? WHERE id = ?", title, content, idPost)
	if err != nil {
		return err
	}

	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id == ?", idPost)
	if err != nil {
		return err
	}

	err = PostModel.AddCategoriesPost(idPost, Categories)
	if err != nil {
		return err
	}

	return nil
}

// add categories for post
func (pm *PostModel) AddCategoriesPost(postID int, categories []string) error {
	ids := []int{}

	for _, categoryName := range categories {
		var categoryID int
		err := pm.DB.QueryRow(
			"SELECT id FROM categories WHERE category_name = ?",
			categoryName,
		).Scan(&categoryID)
		if err != nil {
			return err
		}
		ids = append(ids, categoryID)
	}

	stmt, err := pm.DB.Prepare("INSERT INTO categories_posts (category_id, post_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, categoryID := range ids {
		_, err = stmt.Exec(categoryID, postID)
		if err != nil {
			return err
		}
	}

	return nil
}

// delete post
func (PostModel *PostModel) DeletePost(userId, idPost int) error {
	_, err := PostModel.DB.Exec("DELETE FROM PostTable WHERE id = ? AND user_id = ?", idPost, userId)
	if err != nil {
		return err
	}
	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id = ?", idPost)
	if err != nil {
		return err
	}
	return nil
}
