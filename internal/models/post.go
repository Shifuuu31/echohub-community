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

type CategoriePost struct {
	ID           int
	categorie_id int
	Post_id      int
}

type PostModel struct {
	DB *sql.DB
}

func (PostModel *PostModel) GetCategorys() ([]Categorie, error) {
	Categories := []Categorie{}

	rowsDB, err := PostModel.DB.Query("SELECT id,categorie_name FROM Categories")
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

func (PostModel *PostModel) GetPosts() ([]Post, error) {
	posts := []Post{}

	rowsDB, err := PostModel.DB.Query("SELECT id, user_id, title, post_content, creation_date FROM PostTable ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rowsDB.Close()

	for rowsDB.Next() {
		post := Post{}
		err := rowsDB.Scan(&post.ID, &post.User_id, &post.Title, &post.Post_content, &post.Creation_date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	err = rowsDB.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (PostModel *PostModel) CreatePost(title, content string) (int, error) {
	var id int

	err := PostModel.DB.QueryRow("INSERT INTO PostTable (title, user_id, post_content) VALUES (?, ?, ?) RETURNING id", title, 10, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (PostModel *PostModel) GetIdsCategories(categories []string) ([]int, error) {
	ids := []int{}
	for i := 0; i < len(categories); i++ {
		categorie := Categorie{}

		cmd, err := PostModel.DB.Query("SELECT id, categorie_name FROM Categories WHERE categorie_name = $1", categories[i])
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

func (PostModel *PostModel) AddCategoriePost(post_id int, ids []int) error {
	query := "INSERT INTO Categories_Posts (categorie_id,post_id) VALUES (?,?)"
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

func (PostModel *PostModel) GetUsersNames(posts []Post) ([]string, error) {
	usernames := []string{}

	for _, post := range posts {
		query := "SELECT username FROM UserTable WHERE id = ?"
		var username string
		err := PostModel.DB.QueryRow(query, post.User_id).Scan(&username)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

func (PostModel *PostModel) GetCategoriesNames(posts []Post) ([][]string, error) {
	categoriesNames := make([][]string, len(posts))

	for i, post := range posts {
		queryIDs := "SELECT categorie_id FROM Categories_Posts WHERE post_id = ?"
		rows, err := PostModel.DB.Query(queryIDs, post.ID)
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
			err := PostModel.DB.QueryRow(queryNames, categoryID).Scan(&categoryName)
			if err != nil {
				return nil, err
			}
			categoryNames = append(categoryNames, categoryName)
		}

		categoriesNames[i] = categoryNames
	}

	return categoriesNames, nil
}

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

func (PostModel *PostModel) UpdatePost(idPost int) (string, string, []string, error) {
	post := Post{}

	err := PostModel.DB.QueryRow("SELECT user_id,title,post_content FROM PostTable WHERE id = $1", idPost).Scan(&post.User_id, &post.Title, &post.Post_content)
	if err != nil {
		return "", "", nil, err
	}

	rows, err := PostModel.DB.Query("SELECT categorie_id FROM Categories_Posts WHERE post_id = ?", idPost)
	if err != nil {
		return "", "", nil, err
	}
	defer rows.Close()

	var categories []CategoriePost

	for rows.Next() {
		categoriePost := CategoriePost{
			Post_id: idPost,
		}
		if err := rows.Scan(&categoriePost.categorie_id); err != nil {
			return "", "", nil, err
		}
		categories = append(categories, categoriePost)
	}

	if err = rows.Err(); err != nil {
		return "", "", nil, err
	}

	categorys, err := PostModel.GetCategorys()
	if err != nil {
		return "", "", nil, err
	}

	var selected []string

	for i := 0; i < len(categorys); i++ {
		for j := 0; j < len(categories); j++ {
			if categorys[i].ID == categories[j].categorie_id {
				selected = append(selected, categorys[i].Categorie_name)
			}
		}
	}

	return post.Title, post.Post_content, selected, nil
}

func (PostModel *PostModel) EditPost(idPost int, title string, content string) (err error) {
	_, err = PostModel.DB.Exec("UPDATE PostTable SET title = $1, content = $2 WHERE ID = $3", title, content, idPost)
	if err != nil {
		return err
	}

	return nil
}
