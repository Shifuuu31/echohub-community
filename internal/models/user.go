package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents a user in the system.
type User struct {
	ID             int
	UserName       string
	Email          string
	HashedPassword string
	CreationDate   time.Time
	UserType       string
}

// UserModel handles user-related database operations.
type UserModel struct {
	DB *sql.DB
}

type contextKey string

var (
	UserIDKey   contextKey = "UserID"
	UserTypeKey contextKey = "UserType"
)

func (user *UserModel) RetrieveUser( r *http.Request) (*User, Error){
	var foundUser User
	var err error
	userErr := Error{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		SubMessage: "Unable to retrieve user information.",
		Type: "server",
	}
	// fmt.Println("[",Error{}.Type,"]")

	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		return &User{}, userErr
	}
	userType, ok := r.Context().Value(UserTypeKey).(string)
	if !ok {
		return &User{}, userErr
		
	}
	if userID == 0 {
		return &foundUser, Error{}
	}
	foundUser, err = user.FindUserByID(userID)
	if err != nil {
		return &User{}, userErr
		
	}
	foundUser.UserType = userType
	return &foundUser, Error{}
}

// FindUserByID retrieves a user by their ID from the database.
func (user *UserModel) FindUserByID(userID int) (foundUser User, err error) {
	// foundUser := &User{}

	selectStmt := `SELECT id, username, email FROM UserTable WHERE id = ?`
	err = user.DB.QueryRow(selectStmt, userID).Scan(&foundUser.ID, &foundUser.UserName, &foundUser.Email)
	if err != nil {
		return foundUser, err
	}

	return foundUser, nil
}

// ValidateUserCredentials verifies a user's credentials and returns their ID or errors.
func (user *UserModel) ValidateUserCredentials(username, password string) (UserID int, errors []string) {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" {
		return -1, []string{"Username is required."}
	}
	if password == "" {
		return -1, []string{"Password is required."}
	}

	hashedPassword := ""
	selectStmt := `	SELECT id, hashed_password
					FROM UserTable WHERE username = ?`

	userErr := user.DB.QueryRow(selectStmt, username).Scan(&UserID, &hashedPassword)
	if UserID < 1 || userErr != nil {
		errors = append(errors, "User not found.")
	} else {
		passErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if passErr != nil {
			fmt.Println(passErr.Error())
			errors = append(errors, "Invalid password.")
		}
	}

	if len(errors) != 0 {
		return -1, errors
	}
	
	return UserID, errors
}

// InsertUser adds a new user to the database.
func (user *UserModel) InsertUser(newUser User) (err error) {
	insertStmt := `INSERT INTO UserTable (username, email, hashed_password) VALUES (?, ?, ?)`
	_, err = user.DB.Exec(insertStmt, newUser.UserName, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

type NewUserInfo struct {
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RepeatedPass string `json:"rPassword"`
}
type Response struct {
	Messages []string `json:"messages"`
}

// ValidateNewUser validates new user details and creates a User struct if valid.
func (user *UserModel) ValidateNewUser(new NewUserInfo) (newUser User, errors Response) {
	var err error
	if newUser.UserName, err = user.usernameCheck(new.UserName); err != nil {
		errors.Messages = append(errors.Messages, err.Error())
	}
	if newUser.Email, err = user.emailCheck(new.Email); err != nil {
		errors.Messages = append(errors.Messages, err.Error())
	}

	if new.Password, err = passwordCheck(new.Password, new.RepeatedPass); err != nil {
		errors.Messages = append(errors.Messages, err.Error())
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(new.Password), 12)
	if err != nil {
		errors.Messages = append(errors.Messages, err.Error())
	}
	newUser.HashedPassword = string(hash)

	return newUser, errors
}

// usernameCheck ensures the username is valid and unique.
func (user *UserModel) usernameCheck(username string) (string, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	if username == "" {
		return "", errors.New("Username is required.")
	}
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
	selectStmt := `SELECT username FROM UserTable WHERE username = ?`
	var foundUsername string
	err := user.DB.QueryRow(selectStmt, username).Scan(&foundUsername)
	if err == nil || foundUsername == username {
		return "", errors.New("User '" + username + "' already exists")
	}
	return username, nil
}

// emailCheck ensures the email is valid and unique.
func (user *UserModel) emailCheck(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", errors.New("Email is required.")
	}
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

	selectStmt := `SELECT email FROM UserTable WHERE email = ?`
	var foundEmail string
	err := user.DB.QueryRow(selectStmt, email).Scan(&foundEmail)
	if err == nil || foundEmail == email {
		return "", errors.New("this email '" + email + "' is already registered! Please Log in")
	}
	return email, nil
}

// passwordCheck ensures the password meets security criteria.
func passwordCheck(password, repeatedPassword string) (string, error) {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	if password == "" {
		return "", errors.New("Password is required.")
	}

	if len(password) < 8 || len(password) > 64 {
		return "", errors.New("password must be between 8 and 64 characters")
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
	if repeatedPassword == "" {
		return "", errors.New("Confirm your password.")
	}
	if password != repeatedPassword {
		return "", errors.New("password and repeated must be identical")
	}

	return password, nil
}
