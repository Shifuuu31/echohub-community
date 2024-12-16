package main

import (
	"database/sql"
	"fmt"
	"forum/cmd/web/handlers"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatalln(err)
	}
	handlers.Routes()

	fmt.Println("listening in port : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
