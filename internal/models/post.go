package models

import (
	"database/sql"
	"time"
)

type Post struct {
	ID            int       `json:id`
	User_id       int       `json:user_id`
	Title         string    `json:title`
	Post_content  string    `json:post_content`
	Category_id   int       `json:category_id`
	Creation_date time.Time `json:creation_date`
}

type Categories struct {
	ID            int
	Category_name string
	Post_id       int
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
	defer rowsDB.Close()

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

// categorys []string
func (post *PostModel) CreatePost(title, content string) error {
	query := "INSERT INTO postTable (title,user_id , post_content, category_id) VALUES (?, ?, ?, ?)"
	cmd, err := post.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer cmd.Close()

	_, err = cmd.Exec(title, 1, content, 4)
	if err != nil {
		return err
	}

	return nil
}

func (post *PostModel) GetIdsCategorys(categorys []string) ([]int, error) {
	query := "SELECT id,category_name FROM Categories"
	cmd, err := post.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer cmd.Close()

	category := Categories{}
	ids := []int{}
	for cmd.Next() {
		err = cmd.Scan(&category.ID, &category.Category_name)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(categorys); i++ {
		if categorys[i] == category.Category_name {
			ids = append(ids, category.ID)
		}
	}

	return ids, nil
}
