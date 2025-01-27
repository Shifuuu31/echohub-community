package models

import (
	"database/sql"
)

type Interaction struct {
	EntityID   int    `json:"entityId"`
	EntityType string `json:"entityType"` 
	Liked      bool   `json:"liked"`     
}

type TotalLikesDislikes struct {
	Likes    int  `json:"likes"`
	Dislikes int  `json:"dislikes"`
	IsLiked  sql.NullBool `json:"isLiked"` 
}

type LikesDislikesModel struct {
	DB *sql.DB
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

// TotalLikesDislikes 
func (ldm *LikesDislikesModel) GetTotalLikesDislikes(entityID int, entityType string, userID int) (TotalLikesDislikes, error) {
    var total TotalLikesDislikes

    // Count likes
    query := `
        SELECT COUNT(*) FROM Likes_Dislikes
        WHERE entity_id = ? AND entity_type = ? AND liked = ?
    `
    err := ldm.DB.QueryRow(query, entityID, entityType, true).Scan(&total.Likes)
    if err != nil {
        return total, err
    }

    // Count dislikes
    err = ldm.DB.QueryRow(query, entityID, entityType, false).Scan(&total.Dislikes)
    if err != nil {
        return total, err
    }

	// Check if user liked/disliked
	query = `
        SELECT liked FROM Likes_Dislikes
        WHERE entity_id = ? AND entity_type = ? AND user_id = ?
    `
    err = ldm.DB.QueryRow(query, entityID, entityType, userID).Scan(&total.IsLiked)
    if err != nil && err != sql.ErrNoRows {
        return total, err
    }

    return total, nil
}
