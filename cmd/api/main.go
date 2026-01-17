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

// @title           EchoHub Community Forum API
// @version         1.0
// @description     API documentation for EchoHub Community Forum
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name userSession
// @description Session-based authentication using cookies

func main() {
	models.LoadTemplates()
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
