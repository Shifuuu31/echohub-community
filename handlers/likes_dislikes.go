package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Interactions struct {
	EntityId   int    `json:"entityId"`
	EntityType string `json:"entityType"`
	Liked      bool   `json:"liked"`
}

type Total struct {
	Likes    int  `json:"likes"`
	Dislikes int  `json:"dislikes"`
	IsLiked  bool `json:"isliked"`
}

func LikeDislikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check if user is logged in :

		// userId, ok := r.Context().Value(userIDKey).(int)
		// if !ok || userID == 0 {
		// 	http.Error(w, "Unauthorized: Please log in to like/dislike", http.StatusUnauthorized)
		// 	return
		// }

		var request Interactions
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// remove after adding middelware part above 
		userId := 1

		// check if already liked / disliked
		var alreadyExist bool
		query := `
			SELECT liked FROM Likes_Dislikes
			WHERE entity_id = ? AND user_id = ? AND entity_type = ?
		`
		err := db.QueryRow(query, request.EntityId, userId, request.EntityType).Scan(&alreadyExist)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// if user already interacted with the post
		if err == nil {
			// if user try to like/dislike the same way, remove interaction
			if (alreadyExist && request.Liked) || (!alreadyExist && !request.Liked) {
				_, err := db.Exec(`
					DELETE FROM Likes_Dislikes
					WHERE entity_id = ? AND user_id = ? AND entity_type = ?
				`, request.EntityId, userId, request.EntityType)
				if err != nil {
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				// send response with updated counts
				var total Total
				if err := total.CountTotal(db, request.EntityId, request.EntityType, userId); err != nil {
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(total)
				return
			}

			// ifuser switching from like to dislike or reverse
			_, err := db.Exec(`
				DELETE FROM Likes_Dislikes
				WHERE entity_id = ? AND user_id = ? AND entity_type = ?
			`, request.EntityId, userId, request.EntityType)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		// add new like/dislike
		_, err = db.Exec(`
			INSERT INTO Likes_Dislikes (entity_id, user_id, entity_type, liked)
			VALUES (?, ?, ?, ?)
		`, request.EntityId, userId, request.EntityType, request.Liked)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// updated like/dislike total
		var total Total
		if err := total.CountTotal(db, request.EntityId, request.EntityType, userId); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(total)
	}
}

func (t *Total) CountTotal(db *sql.DB, entityId int, entityType string, userId int) error {
	// likes
	query := `
		SELECT COUNT(*) FROM Likes_Dislikes
		WHERE entity_id = ? AND entity_type = ? AND liked = ?
	`
	err := db.QueryRow(query, entityId, entityType, true).Scan(&t.Likes)
	if err != nil {
		return err
	}

	// dislikes
	err = db.QueryRow(query, entityId, entityType, false).Scan(&t.Dislikes)
	if err != nil {
		return err
	}

	// Check if user already liked/disliked
	var exist bool
	query = `
		SELECT EXISTS (
			SELECT 1 FROM Likes_Dislikes
			WHERE entity_id = ? AND user_id = ? AND entity_type = ?
		)
	`
	err = db.QueryRow(query, entityId, userId, entityType).Scan(&exist)
	if err != nil {
		return err
	}
	t.IsLiked = exist

	return nil
}

/*
user clicks
          |
          v
backend checks cxisting interaction
          |
          v
already user have an interaction?
          |
          v
yup -----------------------> nop
|                            |
v                            v
delete existing one --> add new one
          |                         |
          v                         v
return updated counts <-------------|
          |
          v
front updates ui
*/
