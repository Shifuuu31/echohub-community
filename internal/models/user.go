package models

import (
	"time"
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
