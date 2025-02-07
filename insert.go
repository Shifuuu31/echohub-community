package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	args := os.Args[1:]

	if len(args) == 3 {

		categoryID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalln(err)
		}

		start, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalln(err)
		}

		end, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalln(err)
		}


		for i := start; i <= end; i++ {
			query := "INSERT INTO categories_posts (category_id, post_id) VALUES (?, ?)"

			stmt, err := db.Prepare(query)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(categoryID, i)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Data inserted successfully!")
		}

	}
}
