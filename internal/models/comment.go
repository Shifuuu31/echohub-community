package models

import (
	"database/sql"
	"net/http"
	"time"
)

type Comment struct {
	ID           int
	PostID       int
	UserID       int
	UserName     string
	ProfileImg   string
	Content      string
	LikeCount    int
	DislikeCount int
	Reaction     string
	CreationDate time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (comment *CommentModel) GetPostComments(postID int, userID int) ([]Comment, Error) {
	comments := []Comment{}
	cmd := `
			SELECT CommentTable.id, CommentTable.post_id, CommentTable.user_id, UserTable.username, UserTable.profile_img, CommentTable.comment_content, CommentTable.creation_date 
			FROM CommentTable 
			JOIN UserTable ON CommentTable.user_id = UserTable.id 
			WHERE CommentTable.post_id = ? 
			ORDER BY CommentTable.creation_date DESC`

	rows, err := comment.DB.Query(cmd, postID)
	if err != nil {
		return nil, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}

	defer rows.Close()

	for rows.Next() {
		var Comment Comment
		err := rows.Scan(
			&Comment.ID,
			&Comment.PostID,
			&Comment.UserID,
			&Comment.UserName,
			&Comment.ProfileImg,
			&Comment.Content,
			&Comment.CreationDate,
		)
		if err != nil {
			return nil, Error{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Type:       "server",
			}
		}

		var LDErr Error
		Comment.LikeCount, Comment.DislikeCount, LDErr = GetLikesDislikesCount(comment.DB, Comment.ID, "comment")
		if LDErr.Message != "" {
			return nil, LDErr
		}

		if userID > 0 {
			Reaction, ReactionErr := GetReaction(comment.DB, Comment.ID, "comment", userID)
			if ReactionErr.Message != "" {
				return nil, LDErr
			}
			Comment.Reaction = Reaction
		}

		comments = append(comments, Comment)
	}

	if err = rows.Err(); err != nil {
		return nil, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	return comments, Error{}
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
