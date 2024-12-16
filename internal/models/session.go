package models

import (
	"database/sql"
	"fmt"
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

func (session *SessionModel) GenerateNewSession(userID int) (newSession Session, err error) {
	newToken, err := uuid.NewV4()
	if err != nil {
		return Session{}, err
	}
	newSession = Session{
		UserID:         userID,
		Token:          newToken.String(),
		ExpirationDate: time.Now().Add(7 * 24 * time.Hour), // expire after one week
	}
	return newSession, err
}

func (session *SessionModel) InsertSession(newSession Session) (newCookie http.Cookie, err error) {
	fmt.Println("Inserting Session")

	insertStmt := "INSERT INTO UserSessions (user_id, session_token, expiration_date) VALUES (?, ?, ?);"

	result, err := session.DB.Exec(insertStmt, newSession.UserID, newSession.Token, newSession.ExpirationDate)
	if err != nil {
		return newCookie, err
	}
	fmt.Println("2")
	lastInsertID, err := result.LastInsertId()
	fmt.Println("3")
	if err != nil || lastInsertID != int64(newSession.UserID) {
		return newCookie, err
	}

	newCookie = http.Cookie{
		Name:  "userSession",
		Value: newSession.Token,
		Path:  "/", // valid for the entire site
		// HttpOnly: true,	// Prevent access via JavaScript (mitigates XSS)
		// Secure:   true,	// Transmit only over HTTPS
		Expires: newSession.ExpirationDate,
	}

	fmt.Println("4")
	return newCookie, nil
}

// func (session *SessionModel) ValidateSession(UserID int) (newSession Session, err error) {
// }

// func (session *SessionModel) DeleteSession(UserID int) (newSession Session, err error) {
// }
