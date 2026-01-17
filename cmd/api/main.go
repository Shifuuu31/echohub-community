package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"echohub-community/cmd/web/handlers"
	"echohub-community/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./internal/database/echohub-community.db")
	if err != nil {
		log.Fatalln(err)
	}

	webForum := handlers.WebApp{
		Users: &models.UserModel{
			DB: db,
		},
		Sessions: &models.SessionModel{
			DB: db,
		},
		Posts: &models.PostModel{
			DB: db,
		},
		Comments: &models.CommentModel{
			DB: db,
		},
		LikesDislikes: &models.LikesDislikesModel{
			DB: db,
		},
	}

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}
	go webForum.Sessions.CleanupExpiredSessions()

	server := http.Server{
		Addr:    port,
		Handler: webForum.Router(),
	}

	log.Println("server listening on http://localhost" + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
