package pages

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
)

type ViewPostController struct {
	postRepository *repository.PostRepository
}

func (c *ViewPostController) GetPost(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("post_id") {
		io.WriteString(w, "You should provide 'post_id'!\n")
		return
	}

	postId, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		io.WriteString(w, "'post_id' should be a number!\n")
		return
	}
	post, err := c.postRepository.GetPost(int64(postId))
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, "Post can't be represented!\n")
		return
	}

	io.WriteString(w, fmt.Sprintf("Post %d with post title: %s!\n", post.Id, post.Title))
}

func NewViewPostController(postRepository *repository.PostRepository) *ViewPostController {
	return &ViewPostController{
		postRepository: postRepository,
	}
}
