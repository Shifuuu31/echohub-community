package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/cmd/web/handlers"
	"forum/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var webForum handlers.WebApp
	db, err := sql.Open("sqlite3", "./internal/database/forum.sqlite")
	if err != nil {
		log.Fatalln(err)
	}

	webForum.Post = &models.PostModel{}
	webForum.Post.DB = db

	server := http.Server{
		Addr:    ":5000",
		Handler: webForum.Routes(),
	}

	fmt.Println("listening in port : http://localhost:5000")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
