package models

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	PostId         int
	PostUserName   string
	PostTime       time.Time
	PostTitle      string
	PostContent    string
	PostCategories []string
	LikeCount      int
	DislikeCount   int
	CommentsCount  int
}

type Category struct {
	ID               int
	CategoryName     string
	CategoryIconPath string
}

type categoryPost struct {
	ID          int
	category_id int
	Post_id     int
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
func (postModel *PostModel) GetPosts(startId int, category string) (posts []Post, postsErr Error) {
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
					PostTable.title,
					PostTable.content,
					PostTable.creation_date
				FROM
					PostTable
					JOIN UserTable ON UserTable.id = PostTable.user_id
				WHERE
					UserTable.id = ?
				ORDER BY
					PostTable.id DESC
				LIMIT
					10;`
		args = append(args, startId)

	case "LikedPosts":
		// Placeholder for "LikedPosts" case

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
		err = rows.Scan(&post.PostId, &post.PostUserName, &post.PostTitle, &post.PostContent, &post.PostTime)
		if err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}

		if post.PostCategories, err = postModel.GetPostCategories(post.PostId); err != nil {
			return []Post{}, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}
		post.PostTime = post.PostTime.UTC()
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
func (PostModel *PostModel) GetPostCategories(postId int) (postCategories []string, err error) {
	query := `SELECT
			      Categories.category_name
			  FROM
			      Categories_Posts
			      JOIN Categories ON Categories_Posts.category_id = Categories.id
			  WHERE
			      Categories_Posts.post_id = ?;`

	rows, err := PostModel.DB.Query(query, postId)
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

// type PostData struct {
// 	Id         string   `json:"id`
// 	Title      string   `json:"title"`
// 	Content    string   `json:"content"`
// 	Categories []string `json:"categories"`
// }

// // check form create/update
// func (PostModel *PostModel) CheckErrors(postData PostData, user *User) (postErr Error) {
// 	if len(postData.Categories) == 0 {
// 		return Error{
// 			User:       user,
// 			StatusCode: http.StatusBadRequest,
// 			Type:       "client",
// 			Message:    "Select at least one category",
// 		}
// 	}
// 	if postData.Title == "" {
// 		return Error{
// 			User:       user,
// 			StatusCode: http.StatusBadRequest,
// 			Type:       "client",
// 			Message:    "Title is reqared",
// 		}
// 	}
// 	if postData.Content == "" {
// 		return Error{
// 			User:       user,
// 			StatusCode: http.StatusBadRequest,
// 			Type:       "client",
// 			Message:    "Content is reqared",
// 		}
// 	}
// 	if len(postData.Title) > 70 {
// 		return Error{
// 			User:       user,
// 			StatusCode: http.StatusBadRequest,
// 			Type:       "client",
// 			Message:    "Title length must be up to 70 charactaire",
// 		}
// 	}
// 	if len(postData.Content) > 5000 {
// 		return Error{
// 			User:       user,
// 			StatusCode: http.StatusBadRequest,
// 			Type:       "client",
// 			Message:    "Content length must be up to 70 charactaire",
// 		}
// 	}
// 	return postErr
// }

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
func (PostModel *PostModel) UpdatePost(user_id, idPost int) (post Post, err error) {
	err = PostModel.DB.QueryRow("SELECT p.id,p.title,p.content FROM PostTable p WHERE id = ? AND p.user_id = ?", idPost, user_id).Scan(&post.PostId, &post.PostTitle, &post.PostContent)
	if err != nil {
		return Post{}, errors.New("no post with this ID : ")
	}

	if post.PostCategories, err = PostModel.GetPostCategories(post.PostId); err != nil {
		return Post{}, err
	}

	return post, nil
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

	idsCategoreis, err := PostModel.GetIdsCategories(Categories)
	if err != nil {
		return err
	}

	err = PostModel.AddCategoriesPost(idPost, idsCategoreis)
	if err != nil {
		return err
	}

	return nil
}

// get ids of categories
func (PostModel *PostModel) GetIdsCategories(Categories []string) ([]int, error) {
	ids := []int{}
	for i := 0; i < len(Categories); i++ {
		category := Category{}

		err := PostModel.DB.QueryRow("SELECT id FROM Categories WHERE category_name = ?", Categories[i]).Scan(&category)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}

// add categories for post
func (PostModel *PostModel) AddCategoriesPost(post_id int, ids []int) error {
	query := "INSERT INTO Categories_Posts (category_id,post_id) VALUES (?,?)"
	cmd, err := PostModel.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer cmd.Close()

	for i := 0; i < len(ids); i++ {
		_, err = cmd.Exec(ids[i], post_id)
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
