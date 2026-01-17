package models

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	schemas := []string{
		`CREATE TABLE UserTable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			email TEXT UNIQUE,
			hashed_password TEXT,
			gender TEXT,
			profile_img TEXT,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE Likes_Dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			entity_id INTEGER,
			entity_type TEXT,
			liked BOOLEAN,
			UNIQUE(user_id, entity_id, entity_type)
		);`,
		`CREATE TABLE CommentTable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER,
			user_id INTEGER,
			comment_content TEXT,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE PostTable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			title TEXT,
			content TEXT,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE Categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_name TEXT UNIQUE,
			category_icon_path TEXT
		);`,
		`CREATE TABLE Categories_Posts (
			category_id INTEGER,
			post_id INTEGER,
			PRIMARY KEY (category_id, post_id)
		);`,
	}

	for _, s := range schemas {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("failed to create schema: %v\nQuery: %s", err, s)
		}
	}

	// Verify Tables exist
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		t.Fatalf("Failed to query sqlite_master: %v", err)
	}
	defer rows.Close()
	tables := ""
	for rows.Next() {
		var name string
		rows.Scan(&name)
		tables += name + " "
	}
	t.Logf("Created tables: %s", tables)

	return db
}
