package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID            int
	User_id       int
	Title         string
	Post_content  string
	Category_id   int
	Creation_date time.Time
}

type PostModel struct {
	DB *sql.DB
}

func (post *PostModel) GetPosts() ([]Post, error) {
	posts := []Post{}

	cmd := "SELECT id, user_id, title, post_content, category_id, creation_date FROM PostTable ORDER BY id DESC"
	rowsDB, err := post.DB.Query(cmd)
	if err != nil {
		return nil, err
	}

	for rowsDB.Next() {
		pst := Post{}
		err := rowsDB.Scan(&pst.ID, &pst.User_id, &pst.Title, &pst.Post_content, &pst.Category_id, &pst.Creation_date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, pst)
	}

	err = rowsDB.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (post *PostModel) CreatePost(title, content string, categorys []string) error {
	cmd := "INSERT INTO postTable (title)"

	return nil
}

// func (post *PostModel) GetUser(user_id int) (string, error) {
// }
