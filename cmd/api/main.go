package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"forum/cmd/web"
	"forum/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var webForum web.WebApp

	db, err := sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatalln(err)
	}
	webForum.Users = &models.UserModel{}
	webForum.Users.DB = db
	webForum.Sessions = &models.SessionModel{}
	webForum.Sessions.DB = db

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}

	server := http.Server{
		Addr:    port,
		Handler: webForum.Routes(),
	}

	log.Println("server listening on http://localhost" + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
