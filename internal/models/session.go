package models

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID             int
	UserID         int
	Token          string
	ExpirationDate time.Time
}

type SessionModel struct {
	DB *sql.DB
}

func (session *SessionModel) GenerateNewSession(userID int, remember string) (newSession Session, err error) {
	exp := 24 * time.Hour
	if remember == "on" {
		exp *= 30
	}
	newToken, err := uuid.NewV4()
	if err != nil {
		return Session{}, err
	}
	newSession = Session{
		UserID:         userID,
		Token:          newToken.String(),
		ExpirationDate: time.Now().Add(exp), // expire after 1 day
	}
	return newSession, err
}

func (session *SessionModel) InsertOrUpdateSession(newSession Session) (newCookie http.Cookie, err error) {
	insertOrUpdateStmt := `	INSERT INTO UserSessions (user_id, session_token, expiration_date) VALUES (?, ?, ?)
							ON CONFLICT(user_id) 
							DO UPDATE SET session_token = excluded.session_token, expiration_date = excluded.expiration_date`

	_, err = session.DB.Exec(insertOrUpdateStmt, newSession.UserID, newSession.Token, newSession.ExpirationDate)
	if err != nil {
		return newCookie, err
	}

	newCookie = http.Cookie{
		Name:  "userSession",
		Value: newSession.Token,
		Path:  "/", // valid for the entire site
		// HttpOnly: true,	// Prevent access via JavaScript (mitigates XSS)
		// Secure:   true,	// Transmit only over HTTPS
		Expires: newSession.ExpirationDate,
		// MaxAge: newSession.ExpirationDate.Second(),
	}

	return newCookie, nil
}

func (session *SessionModel) ValidateSession(sessionToken string) (userID int, err error) {
	selectStmt := `	SELECT user_id, expiration_date 
        			FROM UserSessions 
        			WHERE session_token = ?`
	var expirationDate time.Time

	err = session.DB.QueryRow(selectStmt, sessionToken).Scan(&userID, &expirationDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid session token")
		}
		return 0, err
	}

	if time.Now().After(expirationDate) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}

func (session *SessionModel) DeleteSession(sessionToken string) error {
	deleteStmt := `	DELETE FROM UserSessions 
					WHERE session_token = ?`
	result, err := session.DB.Exec(deleteStmt, sessionToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected() // Check if the session was already deleted
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}
