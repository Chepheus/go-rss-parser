package pages

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
)

type RootController struct {
	postRepository *repository.PostRepository
}

func (c *RootController) GetRoot(w http.ResponseWriter, r *http.Request) {
	posts, err := c.postRepository.GetPosts(1, 10)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, "This is my website!\n")
	fmt.Println(posts)
}

func NewRootController(postRepository *repository.PostRepository) *RootController {
	return &RootController{
		postRepository: postRepository,
	}
}
