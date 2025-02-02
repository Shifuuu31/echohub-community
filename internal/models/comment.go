package models

import (
	"database/sql"
	"errors"
	"time"
)

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	UserName     string
	Content      string
	CreationDate time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (comment *CommentModel) Comments(postID int) ([]Comment, error) {
	comments := []Comment{}
	cmd := `
        SELECT c.id, c.post_id, c.user_id, u.username, c.comment_content, c.creation_date 
        FROM CommentTable c 
        JOIN UserTable u ON c.user_id = u.id 
        WHERE c.post_id = ? 
        ORDER BY c.creation_date DESC`

	rows, err := comment.DB.Query(cmd, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.UserName,
			&comment.Content,
			&comment.CreationDate,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	// fmt.Println(comments)
	return comments, nil
}

func (comment *CommentModel) CreateComment(postID, userID int, content string) error {
	if content == "" {
		return errors.New("comment content cannot be empty")
	}
	query := `
        INSERT INTO CommentTable (post_id, user_id, comment_content, creation_date) 
        VALUES (?, ?, ?, ?)`

	stmt, err := comment.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	creationDate := time.Now()
	_, err = stmt.Exec(postID, userID, content, creationDate)
	return err
}

// func (comment *CommentModel) CreateComment(postID, userID int, content string) error {
// 	query := `
//         INSERT INTO CommentTable (post_id, user_id, comment_content, creation_date)
//         VALUES (?, ?, ?, ?)`

// 	stmt, err := PostModel.DB.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	creationDate := time.Now()
// 	_, err = stmt.Exec(postID, userID, content, creationDate)
// 	return err
// }
