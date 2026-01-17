package models

import (
	"testing"
)

func TestCommentModel_GetPostComments(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO UserTable (username, email, profile_img) VALUES ('commenter', 'c@example.com', '/assets/imgs/default.png')")
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	cm := &CommentModel{DB: db}

	err = cm.CreateComment(1, 1, "This is a test comment")
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}

	comments, respErr := cm.GetPostComments(1, 0)
	if respErr.Message != "" {
		t.Fatalf("Failed to get comments: %v", respErr.Message)
	}

	if len(comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(comments))
	}

	if comments[0].Content != "This is a test comment" {
		t.Errorf("Expected content 'This is a test comment', got '%s'", comments[0].Content)
	}
}
