package models

import (
	"database/sql"
	"fmt"
)

type Post struct {
	PostId         int
	PostUserName   string
	PostTime       string
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

func (postModel *PostModel) GetPosts(start, nbr int, category string) ([]Post, error) {
	var query string
	var args []interface{}

	if category == "All" {
		query = `SELECT id, user_id, title, post_content, creation_date
				  FROM PostTable
				  ORDER BY id DESC`
	} else {
		query = `
		SELECT p.id, p.user_id, p.title, p.post_content, p.creation_date
		FROM PostTable p
		JOIN Categories_Posts cp ON p.id = cp.post_id
		JOIN Categories c ON cp.category_id = c.id
		WHERE c.category_name = $1
		ORDER BY p.id DESC`
		args = append(args, category)
	}

	if nbr > 0 {
		if start > 0 {
			query += " LIMIT $2 OFFSET $3"
			args = append(args, nbr, start)
		} else {
			query += " LIMIT $2"
			args = append(args, nbr)
		}
	}

	fmt.Println("Executing query:", query)
	fmt.Println("Executing args:", args)

	rows, err := postModel.DB.Query(query, args[0], args[1], args[2])
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		post := Post{}
		var userID int
		if err := rows.Scan(&post.PostId, &userID, &post.PostTitle, &post.PostContent, &post.PostTime); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		post.PostUserName, err = postModel.GetUsersName(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for user_id %d: %w", userID, err)
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return posts, nil
}

func (PostModel *PostModel) GetUsersName(userId int) (username string, err error) {
	query := "SELECT username FROM UserTable WHERE id = ?"
	if err = PostModel.DB.QueryRow(query, userId).Scan(&username); err != nil {
		return "", err
	}

	return username, nil
}

func (PostModel *PostModel) GetCategoriesPost(postId int) (postCategories []string, err error) {
	query := "SELECT c.category_name FROM Categories_Posts cp JOIN Categories c ON cp.category_id = c.id WHERE cp.post_id = ?"
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

// func (PostModel *PostModel) CreatePost(title, content string) (int, error) {
// 	var id int

// 	err := PostModel.DB.QueryRow("INSERT INTO PostTable (title, user_id, post_content) VALUES (?, ?, ?) RETURNING id", title, 10, content).Scan(&id)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return id, nil
// }

// func (PostModel *PostModel) GetIdsCategories(Categories []string) ([]int, error) {
// 	ids := []int{}
// 	for i := 0; i < len(Categories); i++ {
// 		category := Category{}

// 		cmd, err := PostModel.DB.Query("SELECT id, category_name FROM Categories WHERE category_name = $1", Categories[i])
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer cmd.Close()

// 		for cmd.Next() {
// 			err = cmd.Scan(&category.ID, &category.CategoryName)
// 			if err != nil {
// 				return nil, err
// 			}
// 			ids = append(ids, category.ID)
// 		}
// 	}

// 	return ids, nil
// }

// func (PostModel *PostModel) AddcategoryPost(post_id int, ids []int) error {
// 	query := "INSERT INTO Categories_Posts (category_id,post_id) VALUES (?,?)"
// 	cmd, err := PostModel.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer cmd.Close()

// 	for i := 0; i < len(ids); i++ {
// 		_, err = cmd.Exec(ids[i], post_id)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (PostModel *PostModel) DeletePost(idPost int) error {
// 	_, err := PostModel.DB.Exec("DELETE FROM PostTable WHERE ID = $1", idPost)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id = $1", idPost)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (PostModel *PostModel) UpdatePost(idPost int) (string, string, []string, error) {
// 	post := Post{}

// 	err := PostModel.DB.QueryRow("SELECT user_id,title,post_content FROM PostTable WHERE id = $1", idPost).Scan(&post.PostId, &post.PostTitle, &post.PostContent)
// 	if err != nil {
// 		return "", "", nil, err
// 	}

// 	rows, err := PostModel.DB.Query("SELECT category_id FROM Categories_Posts WHERE post_id = ?", idPost)
// 	if err != nil {
// 		return "", "", nil, err
// 	}
// 	defer rows.Close()

// 	var Categories []categoryPost

// 	for rows.Next() {
// 		categoryPost := categoryPost{
// 			Post_id: idPost,
// 		}
// 		if err := rows.Scan(&categoryPost.category_id); err != nil {
// 			return "", "", nil, err
// 		}
// 		Categories = append(Categories, categoryPost)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return "", "", nil, err
// 	}

// 	categorys, err := PostModel.GetCategories()
// 	if err != nil {
// 		return "", "", nil, err
// 	}

// 	var selected []string

// 	for i := 0; i < len(categorys); i++ {
// 		for j := 0; j < len(Categories); j++ {
// 			if categorys[i].ID == Categories[j].category_id {
// 				selected = append(selected, categorys[i].CategoryName)
// 			}
// 		}
// 	}

// 	return post.PostTitle, post.PostContent, selected, nil
// }

// func (PostModel *PostModel) EditPost(idPost int, title string, content string, Categories []string) (err error) {
// 	_, err = PostModel.DB.Exec("UPDATE PostTable SET title = $1, post_content = $2 WHERE ID = $3", title, content, idPost)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = PostModel.DB.Exec("DELETE FROM Categories_Posts WHERE post_id == $1", idPost)
// 	if err != nil {
// 		return err
// 	}

// 	idsCategoreis, err := PostModel.GetIdsCategories(Categories)
// 	if err != nil {
// 		return err
// 	}

// 	err = PostModel.AddcategoryPost(idPost, idsCategoreis)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
