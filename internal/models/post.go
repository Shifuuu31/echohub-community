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
	Creation_date time.Time
}

type Categorie struct {
	ID             int
	Categorie_name string
}

type PostModel struct {
	DB *sql.DB
}

func (post *PostModel) GetCategorys() ([]Categorie, error) {
	Categories := []Categorie{}

	cmd := "SELECT id,categorie_name FROM Categories"
	rowsDB, err := post.DB.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		Categorie := Categorie{}
		err := rowsDB.Scan(&Categorie.ID, &Categorie.Categorie_name)
		if err != nil {
			return nil, err
		}
		Categories = append(Categories, Categorie)
	}

	err = rowsDB.Err()
	if err != nil {
		return nil, err
	}

	return Categories, nil
}

func (post *PostModel) GetPosts() ([]Post, error) {
	posts := []Post{}

	cmd := "SELECT id, user_id, title, post_content, creation_date FROM PostTable ORDER BY id DESC"
	rowsDB, err := post.DB.Query(cmd)
	if err != nil {
		return nil, err
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		pst := Post{}
		err := rowsDB.Scan(&pst.ID, &pst.User_id, &pst.Title, &pst.Post_content, &pst.Creation_date)
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

func (post *PostModel) GetIdsCategories(categories []string) ([]int, error) {
	ids := []int{}
	for i := 0; i < len(categories); i++ {
		categorie := Categorie{}
		query := "SELECT id, categorie_name FROM Categories WHERE categorie_name = $1"

		cmd, err := post.DB.Query(query, categories[i])
		if err != nil {
			return nil, err
		}
		defer cmd.Close()

		for cmd.Next() {
			err = cmd.Scan(&categorie.ID, &categorie.Categorie_name)
			if err != nil {
				return nil, err
			}
			ids = append(ids, categorie.ID)
		}
	}

	return ids, nil
}

func (post *PostModel) AddCategoriePost(post_id int, ids []int) error {
	query := "INSERT INTO Categories_Posts (categorie_id,post_id) VALUES (?,?)"
	cmd, err := post.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer cmd.Close()

	for i := 0; i < len(ids); i++ {
		_, err = cmd.Exec(ids[0], post_id)
		if err != nil {
			return err
		}
	}

	return nil
}
