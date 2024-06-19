package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
	"github.com/Chepheus/go-rss-parser/web-server/vo"
	"github.com/Chepheus/go-rss-parser/web-server/web/services"
)

type ViewPostController struct {
	htmlRenderer   services.HTMLRenderer
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

	data := &services.TemplateData{
		"post": vo.NewPostDetails(post.Id, post.Title, post.MainImage.String, post.Content.String, post.PubDate),
	}

	c.htmlRenderer.Render(w, "post", data)
}

func NewViewPostController(htmlRenderer services.HTMLRenderer, postRepository *repository.PostRepository) *ViewPostController {
	return &ViewPostController{
		htmlRenderer:   htmlRenderer,
		postRepository: postRepository,
	}
}
