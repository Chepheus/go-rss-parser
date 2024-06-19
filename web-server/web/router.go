package web

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
	"github.com/Chepheus/go-rss-parser/web-server/web/controllers"
	"github.com/Chepheus/go-rss-parser/web-server/web/services"
)

func StartWebServer(port string, db *sql.DB) {
	htmlRenderer := services.NewHTMLRenderer("./web/templates/")
	postRepository := repository.NewPostRepository(db)
	rootController := controllers.NewRootController(htmlRenderer, postRepository)
	viewPostController := controllers.NewViewPostController(htmlRenderer, postRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootController.GetRoot)
	mux.HandleFunc("/post", viewPostController.GetPost)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
