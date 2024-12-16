package sqlite

import (
	"database/sql"
	"forum/internal/models"
)

type PostModel struct {
	DB *sql.DB
}

func GetPosts(model *PostModel) []models.Post {
	cmd := "SELECT id, user_id, title, post_content, category_id, creation_date FROM PosteTable ORDER BY id DESC"
	rowsDB, err := model.DB.Query(cmd)
	if err != nil {

	}
}
