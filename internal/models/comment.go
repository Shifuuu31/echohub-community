package models

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	UserName     string
	ProfileImg     string
	Content      string
	LikeCount      int
	DislikeCount   int
	CreationDate time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (comment *CommentModel) Comments(postID int) ([]Comment, error) {
	comments := []Comment{}
	cmd := `
        SELECT CommentTable.id, CommentTable.post_id, CommentTable.user_id, UserTable.username, UserTable.profile_img, CommentTable.comment_content, CommentTable.creation_date 
        FROM CommentTable 
        JOIN UserTable ON CommentTable.user_id = UserTable.id 
        WHERE CommentTable.post_id = ? 
        ORDER BY CommentTable.creation_date DESC`

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
			&comment.ProfileImg,
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
	return comments, nil
}

func (comment *CommentModel) CreateComment(postID, userID int, content string) error {
	insertStmt := `
    INSERT INTO CommentTable (post_id, user_id, comment_content) 
    VALUES (?, ?, ?) `

	_, err := comment.DB.Exec(insertStmt, postID, userID, content)
	if err != nil {
		return err
	}

	return nil
}

