package models

import (
	"database/sql"
	"errors"
	"fmt"
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
func (PostModel *PostModel) GetCategories() (Categories []Category, err error) {
	rowsDB, err := PostModel.DB.Query("SELECT id,category_name,category_icon_path FROM Categories")
	if err != nil {
		return nil, err
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		category := Category{}
		err := rowsDB.Scan(&category.ID, &category.CategoryName, &category.CategoryIconPath)
		if err != nil {
			return nil, err
		}
		Categories = append(Categories, category)
	}

	err = rowsDB.Err()
	if err != nil {
		return nil, err
	}

	return Categories, nil
}

// Get maxID of posts
func (postModel *PostModel) GetMaxID() (maxID int, err error) {
	if err = postModel.DB.QueryRow("SELECT p.id	FROM PostTable p ORDER BY p.id DESC LIMIT 1").Scan(&maxID); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return -1, fmt.Errorf("error scanning max ID: %w", err)
	}
	return maxID, nil
}

// get posts from DB with cateogry
func (postModel *PostModel) GetPosts(current int, category string) (posts Post, err error) {
	var query string
	var categoryId int

	var args []interface{}
	if current <= 0 {
		return Post{}, errors.New("no argements")
	}

	switch category {
	case "All":
		query = `SELECT p.id,u.username,p.title,p.post_content,p.creation_date
		FROM PostTable p
		JOIN UserTable u ON u.id = p.user_id
		WHERE p.id = $1;`
		args = append(args, current)
	case "MyPosts":
		query = `SELECT p.id,u.username,p.title,p.post_content,p.creation_date
		FROM PostTable p
		JOIN UserTable u ON u.id = p.user_id
		WHERE user_id = $2;`
		// need user ID
	case "LikedPosts":
		// to be implemented
	default:
		if categoryId, err = strconv.Atoi(category); err != nil {
			return Post{}, errors.New("This category not defined")
		}
		query = `SELECT p.id,u.username,p.title,p.post_content,p.creation_date
		FROM PostTable p
		JOIN Categories_Posts cp ON p.id = cp.post_id
    	JOIN UserTable u ON u.id = p.user_id
		WHERE cp.category_id = $1 AND p.id = $2
		ORDER  BY p.id DESC;`
		args = append(args, categoryId, current)
	}

	stmt, err := postModel.DB.Prepare(query)
	if err != nil {
		return Post{}, err
	}
	defer stmt.Close()

	post := Post{}
	err = stmt.QueryRow(args...).Scan(&post.PostId, &post.PostUserName, &post.PostTitle, &post.PostContent, &post.PostTime)
	if err != nil {
		return Post{}, nil
	}

	if post.PostCategories, err = postModel.GetCategoriesPost(post.PostId); err != nil {
		return Post{}, err
	}
	post.PostTime = post.PostTime.UTC()

	return post, nil
}

// get cateogries of post
func (PostModel *PostModel) GetCategoriesPost(postId int) (postCategories []string, err error) {
	query := `SELECT c.category_name FROM Categories_Posts cp
	JOIN Categories c ON cp.category_id = c.id WHERE cp.post_id = ?`
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

// create post (insert in DB)
func (PostModel *PostModel) CreatePost(title, content string) (int, error) {
	var id int
	err := PostModel.DB.QueryRow("INSERT INTO PostTable (title, user_id, post_content) VALUES (?, ?, ?) RETURNING id", title, 1, content).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// get post by id to update
func (PostModel *PostModel) UpdatePost(idPost int) (post Post, err error) {
	err = PostModel.DB.QueryRow("SELECT p.id,p.title,p.post_content FROM PostTable p WHERE id = $1", idPost).Scan(&post.PostId, &post.PostTitle, &post.PostContent)
	if err != nil {
		return Post{}, errors.New("no post with this ID : ")
	}

	if post.PostCategories, err = PostModel.GetCategoriesPost(post.PostId); err != nil {
		return Post{}, err
	}

	return post, nil
}

// update post
func (PostModel *PostModel) EditPost(idPost int, title string, content string, Categories []string) (err error) {
	_, err = PostModel.DB.Exec("UPDATE PostTable SET title = $1, post_content = $2 WHERE ID = $3", title, content, idPost)
	if err != nil {
		return err
	}

	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id == $1", idPost)
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

		cmd, err := PostModel.DB.Query("SELECT id, category_name FROM Categories WHERE category_name = $1", Categories[i])
		if err != nil {
			return nil, err
		}
		defer cmd.Close()

		for cmd.Next() {
			err = cmd.Scan(&category.ID, &category.CategoryName)
			if err != nil {
				return nil, err
			}
			ids = append(ids, category.ID)
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
func (PostModel *PostModel) DeletePost(idPost int) error {
	_, err := PostModel.DB.Exec("DELETE FROM PostTable WHERE ID = $1", idPost)
	if err != nil {
		return err
	}
	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id = $1", idPost)
	if err != nil {
		return err
	}
	return nil
}
