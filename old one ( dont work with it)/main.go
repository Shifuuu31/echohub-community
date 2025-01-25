package main

import (
	"database/sql"
	"fmt"
	"forum/handlers"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS UserTable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS PostTable (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			post_content TEXT NOT NULL,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS Likes_Dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			entity_id INTEGER NOT NULL,
			entity_type TEXT NOT NULL CHECK(entity_type IN ('post', 'comment')),
			liked BOOLEAN NOT NULL,
			FOREIGN KEY(user_id) REFERENCES UserTable(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

func main() {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Initialize the database
	if err := initDB(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/like-dislike", handlers.LikeDislikeHandler(db))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index2.html")
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
