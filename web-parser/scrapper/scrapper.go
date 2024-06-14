package scrapper

import (
	"github.com/gocolly/colly/v2"
)

type WebPost struct {
	ExternalPostLink string
	MainImageSrc     string
	TextBlocks       []string
}

type WebScrapper struct {
	c *colly.Collector
}

func (ws *WebScrapper) Visit(postUrl string) (*WebPost, error) {
	webPost := WebPost{ExternalPostLink: postUrl}
	ws.c.OnHTML("#main-content > article", func(e *colly.HTMLElement) {
		webPost.MainImageSrc = e.ChildAttr("img[loading=\"eager\"]", "src")

		e.ForEach("div[data-component=\"text-block\"]", func(_ int, e *colly.HTMLElement) {
			webPost.TextBlocks = append(webPost.TextBlocks, e.ChildTexts("p")...)
		})
	})

	return &webPost, ws.c.Visit(postUrl)
}

func NewWebScrapper() *WebScrapper {
	return &WebScrapper{
		c: colly.NewCollector(),
	}
}
