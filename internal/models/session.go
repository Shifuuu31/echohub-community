package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Session represents a user session in the database.
type Session struct {
	ID             int
	UserID         int
	Token          string
	ExpirationDate time.Time
}

// SessionModel handles operations related to user sessions.
type SessionModel struct {
	DB *sql.DB
}

// GenerateNewSession creates a new session for the given user ID.
// If "remember" is true, the session lasts for approximately 30 days; otherwise, it expires in 24 hours.
func (session *SessionModel) GenerateNewSession(userID int, remember bool) (newSession Session, err error) {
	exp := 24 * time.Hour
	if remember {
		exp *= 30 // approximately 30 days
	}

	newToken, err := uuid.NewV4()
	if err != nil {
		return Session{}, err
	}

	newSession = Session{
		UserID:         userID,
		Token:          newToken.String(),
		ExpirationDate: time.Now().Add(exp),
	}

	fmt.Println("\x1b[1;31m", newSession.ExpirationDate, "\x1b[1;39m")
	return newSession, err
}

// InsertOrUpdateSession adds or updates a session in the database.
// It returns an HTTP cookie representing the session.
func (session *SessionModel) InsertOrUpdateSession(newSession Session) (newCookie http.Cookie, err error) {
	insertOrUpdateStmt := `
		INSERT INTO UserSessions (user_id, session_token, expiration_date) VALUES (?, ?, ?)
		
		ON CONFLICT(user_id) 
		DO UPDATE SET session_token = excluded.session_token, expiration_date = excluded.expiration_date`

	_, err = session.DB.Exec(insertOrUpdateStmt, newSession.UserID, newSession.Token, newSession.ExpirationDate)
	if err != nil {
		log.Println(err)
		return newCookie, err
	}

	newCookie = http.Cookie{
		Name:     "userSession",
		Value:    newSession.Token,
		Path:     "/",
		HttpOnly: true,
		Expires:  newSession.ExpirationDate.Add(time.Hour), // accounts for time zone difference
	}

	return newCookie, nil
}

// ValidateSession checks if a session token is valid and returns the associated user ID.
func (session *SessionModel) ValidateSession(sessionToken string) (userID int, sessionErr Error) {
	selectStmt := `
		SELECT user_id, expiration_date 
		FROM UserSessions 
		WHERE session_token = ?`

	var expirationDate time.Time
	if err := session.DB.QueryRow(selectStmt, sessionToken).Scan(&userID, &expirationDate); err != nil {
		if err == sql.ErrNoRows {
			return 0, sessionErr
		}
		return 0, Error{
			// User: &User{},
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
			Type:       "server",
		}
	}

	if time.Now().After(expirationDate) {
		return 0, sessionErr
	}
	return userID, sessionErr
}

// DeleteSession removes a session based on its token.
func (session *SessionModel) DeleteSession(sessionToken string) error {
	deleteStmt := `
		DELETE FROM UserSessions 
		WHERE session_token = ?`

	result, err := session.DB.Exec(deleteStmt, sessionToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}

// deleteExpiredSessions removes all sessions that have expired from the database.
func (session *SessionModel) deleteExpiredSessions() error {
	deleteStmt := `DELETE FROM UserSessions WHERE expiration_date < CURRENT_TIMESTAMP`
	result, err := session.DB.Exec(deleteStmt)
	if err != nil {
		return fmt.Errorf("no Expired session to be deleted: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to fetch rows affected: %v", err)
	}

	log.Printf("Cleanup: Deleted %d expired session(s).", rowsAffected)
	return nil
}

// CleanupExpiredSessions periodically cleans up expired sessions from the database.
// This runs in an infinite loop with a 30-second delay between each execution.
func (session *SessionModel) CleanupExpiredSessions() {
	for {
		time.Sleep(30 * time.Second)
		err := session.deleteExpiredSessions()
		if err != nil {
			log.Printf("Error removing expired sessions: %v", err)
		}
	}
}
