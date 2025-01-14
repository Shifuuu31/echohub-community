package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	UserName       string
	Email          string
	HashedPassword string
	CreationDate   time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (user *UserModel) FindUserByID(userID int) (foundUser User, err error) {
    selectStmt := `	SELECT id, username, email, hashed_password, creation_date
        			FROM UserTable
        			WHERE id = ?`
    err = user.DB.QueryRow(selectStmt, userID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.Email, &foundUser.HashedPassword, &foundUser.CreationDate)

    if err != nil {
        if err == sql.ErrNoRows {
            return foundUser, errors.New("user not found")
        }
        return foundUser, err
    }
	
    return foundUser, nil
}

func (user *UserModel) ValidateUserCreadentials(username, password string) (UserID int, err error) {
	username = strings.TrimSpace(username)
	hashedPassword := ""
	selectStmt := `	SELECT id, username, hashed_password
					FROM UserTable WHERE username = ?`
	err = user.DB.QueryRow(selectStmt, username).Scan(&UserID, &username, &hashedPassword)
	if err != nil {
		return -1, errors.New("wrong username")
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return -1, errors.New("wrong password")
	}

	fmt.Println(UserID, username, password, hashedPassword)
	return UserID, nil
}

func (user *UserModel) InsertUser(newUser User) (err error) {
	insertStmt := `INSERT INTO UserTable (username, email, hashed_password) VALUES (?, ?, ?)`
	_, err = user.DB.Exec(insertStmt, newUser.UserName, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (user *UserModel) ValidateNewUser(username, email, password, repeatedPassword string) (newUser User, err error) {
	if newUser.UserName, err = user.usernameCheck(username); err != nil {
		return User{}, err
	}
	if newUser.Email, err = user.emailCheck(email); err != nil {
		return User{}, err
	}

	if password, err = passwordCheck(password, repeatedPassword); err != nil {
		return User{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return User{}, err
	}
	newUser.HashedPassword = string(hash)

	return newUser, err
}

func (user *UserModel) usernameCheck(username string) (string, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return "", errors.New("username must be between 3 and 20 characters")
	}
	if username[0] == '_' || username[len(username)-1] == '_' {
		return "", errors.New("username cannot start or end with '_'")
	}
	for _, char := range username {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') && !(char >= '0' && char <= '9') && char != '_' {
			return "", errors.New("username can only contain letters(A-Za-z), numbers(0-9) and underscores(_)")
		}
	}
	selectStmt := `	SELECT COUNT(*)
					FROM UserTable WHERE username = ?`
	var count int
	err := user.DB.QueryRow(selectStmt, username).Scan(&count)
	if err != nil || count > 0 {
		return "", err
	}
	return username, nil
}

func (user *UserModel) emailCheck(email string) (string, error) {
	email = strings.TrimSpace(email)

	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return "", errors.New("email must contain '@'")
	}
	if strings.Count(email, "@") > 1 {
		return "", errors.New("email must contain only one '@'")
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:]

	if len(localPart) == 0 {
		return "", errors.New("email must have a local part before '@'")
	}
	if len(domainPart) == 0 {
		return "", errors.New("email must have a domain part after '@'")
	}

	if !strings.Contains(domainPart, ".") {
		return "", errors.New("email domain must contain '.'")
	}

	if domainPart[0] == '.' || domainPart[len(domainPart)-1] == '.' {
		return "", errors.New("email domain cannot start or end with '.'")
	}

	var count int
	selectStmt := `SELECT COUNT(*) FROM UserTable WHERE username = ?`
	err := user.DB.QueryRow(selectStmt, email).Scan(&count)
	if err != nil || count != 0 {
		return "", errors.New("this email '" + email + "' is alreaady registered! Please Log in")
	}
	return email, nil
}

func passwordCheck(password, repeatedPassword string) (string, error) {
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	if len(password) < 8 || len(password) > 64 {
		return "", errors.New("password must be between 8 and 64 characters")
	}

	if password != repeatedPassword {
		return "", errors.New("password and repeated must be identical")
	}

	specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/"

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		case char == ' ':
			return "", errors.New("password cannot contain spaces")
		}
	}

	switch false {
	case hasUpper:
		return "", errors.New("password must contain at least one uppercase letter")
	case hasLower:
		return "", errors.New("password must contain at least one lowercase letter")
	case hasDigit:
		return "", errors.New("password must contain at least one number")
	case hasSpecial:
		return "", errors.New("password must contain at least one special character")
	}

	return password, nil
}
