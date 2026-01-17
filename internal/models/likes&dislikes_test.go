package models

import (
	"testing"
)

func TestLikesDislikesModel_LikeDislike(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Close()

	ldm := &LikesDislikesModel{DB: db}

	// Test liking a post
	err := ldm.LikeDislike(1, "post", 1, true)
	if err != nil {
		t.Fatalf("Failed to like post: %v", err)
	}

	likeCount, dislikeCount, respErr := GetLikesDislikesCount(db, 1, "post")
	if respErr.Message != "" {
		t.Fatalf("Failed to get count: %v", respErr.Message)
	}

	if likeCount != 1 || dislikeCount != 0 {
		t.Errorf("Expected 1 like, 0 dislikes. Got L:%d, D:%d", likeCount, dislikeCount)
	}

	// Toggle to dislike
	err = ldm.LikeDislike(1, "post", 1, false)
	if err != nil {
		t.Fatalf("Failed to dislike post: %v", err)
	}

	likeCount, dislikeCount, _ = GetLikesDislikesCount(db, 1, "post")
	if likeCount != 0 || dislikeCount != 1 {
		t.Errorf("Expected 0 likes, 1 dislike. Got L:%d, D:%d", likeCount, dislikeCount)
	}

	// Remove (third click on same reaction usually removes it in some implementations,
	// but let's see how LikeDislike is implemented)
	err = ldm.LikeDislike(1, "post", 1, false)
	if err != nil {
		t.Fatalf("Failed to remove dislike: %v", err)
	}

	likeCount, dislikeCount, _ = GetLikesDislikesCount(db, 1, "post")
	if likeCount != 0 || dislikeCount != 0 {
		t.Errorf("Expected 0 likes, 0 dislikes after removal. Got L:%d, D:%d", likeCount, dislikeCount)
	}
}
