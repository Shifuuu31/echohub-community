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

// type user struct {
// 	ID       int
// 	username string
// }

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
	query := "INSERT INTO PostTable (title,user_id , post_content) VALUES (?, ?, ?)"
	cmd, err := post.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer cmd.Close()

	_, err = cmd.Exec(title, 1, content)
	if err != nil {
		return err
	}

	return nil
}

func (post *PostModel) GetLastPoID() (int, error) {
	posts := []Post{}

	cmd := "SELECT id, user_id, title, post_content, creation_date FROM PostTable ORDER BY id DESC"
	rowsDB, err := post.DB.Query(cmd)
	if err != nil {
		return 0, err
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		pst := Post{}
		err := rowsDB.Scan(&pst.ID, &pst.User_id, &pst.Title, &pst.Post_content, &pst.Creation_date)
		if err != nil {
			return 0, err
		}
		posts = append(posts, pst)
	}

	err = rowsDB.Err()
	if err != nil {
		return 0, err
	}

	return posts[0].ID, nil
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
		_, err = cmd.Exec(ids[i], post_id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (postModel *PostModel) GetUsersNames(posts []Post) ([]string, error) {
	usernames := []string{}

	for _, post := range posts {
		query := "SELECT username FROM UserTable WHERE id = ?"
		var username string
		err := postModel.DB.QueryRow(query, post.User_id).Scan(&username)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func (postModel *PostModel) GetCategoriesNames(posts []Post) ([][]string, error) {
	categoriesNames := make([][]string, len(posts))

	for i, post := range posts {
		queryIDs := "SELECT categorie_id FROM Categories_Posts WHERE post_id = ?"
		rows, err := postModel.DB.Query(queryIDs, post.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var categoryIDs []int
		for rows.Next() {
			var categoryID int
			if err := rows.Scan(&categoryID); err != nil {
				return nil, err
			}
			categoryIDs = append(categoryIDs, categoryID)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		var categoryNames []string
		for _, categoryID := range categoryIDs {
			queryNames := "SELECT categorie_name FROM Categories WHERE id = ?"
			var categoryName string
			err := postModel.DB.QueryRow(queryNames, categoryID).Scan(&categoryName)
			if err != nil {
				return nil, err
			}
			categoryNames = append(categoryNames, categoryName)
		}

		categoriesNames[i] = categoryNames
	}

	return categoriesNames, nil
}
