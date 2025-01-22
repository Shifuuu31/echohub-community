package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"forum/cmd/web/handlers"
	"forum/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var webForum handlers.WebApp
	db, err := sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatalln(err)
	}
	webForum.Post = &models.PostModel{}
	webForum.Post.DB = db
	port := os.Args[1]
	server := http.Server{
		Addr:    ":"+port,
		Handler: webForum.Routes(),
	}

	fmt.Println("listening in port : http://localhost:"+port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
