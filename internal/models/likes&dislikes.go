package models
import (
	"database/sql"
	"net/http"
)
type Interaction struct {
	EntityID   int    `json:"entityId"`
	EntityType string `json:"entityType"`
	Liked      bool   `json:"liked"`
}
type LikesDislikesModel struct {
	DB *sql.DB
}
func GetLikesDislikesCount(DB *sql.DB, entityID int, entityType string) (likeCount, dislikeCount int, LDErr Error) {
	query := `
	SELECT COUNT(*) FROM Likes_Dislikes
	WHERE entity_id = ? AND entity_type = ? AND liked = ?
    `
	// Count likes
	err := DB.QueryRow(query, entityID, entityType, true).Scan(&likeCount)
	if err != nil {
		return likeCount, dislikeCount, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	// Count dislikes
	err = DB.QueryRow(query, entityID, entityType, false).Scan(&dislikeCount)
	if err != nil {
		return likeCount, dislikeCount, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	return likeCount, dislikeCount, LDErr
}
func GetReaction(DB *sql.DB, entityID int, entityType string, userID int) (state string, ReactionErr Error) {
	var liked bool
	// Check if user liked/disliked
	query := `
        SELECT liked FROM Likes_Dislikes
        WHERE entity_id = ? AND entity_type = ? AND user_id = ?
    `
	err := DB.QueryRow(query, entityID, entityType, userID).Scan(&liked)
	if err != nil {
		if err == sql.ErrNoRows {
			return "none", ReactionErr
		}
		return state, Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}
	if liked == true {
		return "liked", ReactionErr
	}
	return "disliked", ReactionErr
}
// like/dislike a post or comment
func (ldm *LikesDislikesModel) LikeDislike(entityID int, entityType string, userID int, liked bool) error {
	// check if already liked / disliked
	var existingLike bool
	query := `
		SELECT liked FROM Likes_Dislikes
		WHERE entity_id = ? AND user_id = ? AND entity_type = ?
	`
	err := ldm.DB.QueryRow(query, entityID, userID, entityType).Scan(&existingLike)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil {
		// if user try to like/dislike the same way, remove interaction
		if (existingLike && liked) || (!existingLike && !liked) {
			_, err := ldm.DB.Exec(`
				DELETE FROM Likes_Dislikes
				WHERE entity_id = ? AND user_id = ? AND entity_type = ?
			`, entityID, userID, entityType)
			return err
		}
		// ifuser switching from like to dislike or reverse
		_, err := ldm.DB.Exec(`
			DELETE FROM Likes_Dislikes
			WHERE entity_id = ? AND user_id = ? AND entity_type = ?
		`, entityID, userID, entityType)
		if err != nil {
			return err
		}
	}
	// Add the new like/dislike
	_, err = ldm.DB.Exec(`
		INSERT INTO Likes_Dislikes (entity_id, user_id, entity_type, liked)
		VALUES (?, ?, ?, ?)
	`, entityID, userID, entityType, liked)
	return err
}