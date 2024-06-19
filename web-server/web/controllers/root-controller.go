package controllers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Chepheus/go-rss-parser/web-server/storage/repository"
	"github.com/Chepheus/go-rss-parser/web-server/vo"
	"github.com/Chepheus/go-rss-parser/web-server/web/services"
)

type RootController struct {
	htmlRenderer   services.HTMLRenderer
	postRepository *repository.PostRepository
}

func (c *RootController) GetRoot(w http.ResponseWriter, r *http.Request) {
	page, err := c.getQueryParam(r, "page", 1)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	limit, err := c.getQueryParam(r, "limit", 10)
	if err != nil {
		io.WriteString(w, err.Error())
	}

	posts, err := c.postRepository.GetPosts(page, limit)
	if err != nil {
		log.Fatal(err)
	}

	var postList []vo.PostListElem
	for _, p := range posts {
		postListElem := vo.NewPostListElem(p.Id, p.Title, p.Thumbnail.String, p.ShortDescription.String, p.PubDate)
		postList = append(postList, postListElem)
	}

	data := &services.TemplateData{
		"title": "Posts List",
		"posts": postList,
	}

	c.htmlRenderer.Render(w, "home", data)
}

func (c *RootController) getQueryParam(r *http.Request, param string, defVal int) (int32, error) {
	page := defVal
	var err error
	if r.URL.Query().Has(param) {
		page, err = strconv.Atoi(r.URL.Query().Get(param))
		if err != nil {
			return int32(defVal), err
		}
	}

	return int32(page), nil
}

func NewRootController(htmlRenderer services.HTMLRenderer, postRepository *repository.PostRepository) *RootController {
	return &RootController{
		htmlRenderer:   htmlRenderer,
		postRepository: postRepository,
	}
}
