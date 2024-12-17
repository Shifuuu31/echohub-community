package sqlite

import (
	"database/sql"
	"time"
)

type Post struct {
	ID            int
	user_id       int
	title         string
	post_content  string
	category_id   int
	creation_date time.Time
}

type PostModel struct {
	DB *sql.DB
}

func GetPosts(model *PostModel) []Post {
	// cmd := "SELECT id, user_id, title, post_content, category_id, creation_date FROM PosteTable ORDER BY id DESC"
	// rowsDB, err := model.DB.Query(cmd)
	// if err != nil {
	// }
	return []Post{}
}
