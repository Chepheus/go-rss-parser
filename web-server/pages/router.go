package pages

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
)

func StartWebServer(port string, db *sql.DB) {
	postRepository := repository.NewPostRepository(db)
	rootController := NewRootController(postRepository)
	viewPostController := NewViewPostController(postRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootController.GetRoot)
	mux.HandleFunc("/post", viewPostController.GetPost)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
