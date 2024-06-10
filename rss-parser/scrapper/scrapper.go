package scrapper

import (
	"log"

	"github.com/gocolly/colly/v2"
)

type PostData struct {
	Title       string
	Description string
	Link        string
	Thumbnail   string
	PubDate     string
}

type RssScrapper struct {
	c *colly.Collector
}

func (rs *RssScrapper) VisitRss(rssUrl string) ([]PostData, error) {
	posts := []PostData{}

	rs.c.OnXML("//channel", func(e *colly.XMLElement) {
		titles := e.ChildTexts("/item/title")
		descriptions := e.ChildTexts("/item/description")
		links := e.ChildTexts("/item/link")
		thumbnails := e.ChildAttrs("/item/media:thumbnail", "url")
		pubDates := e.ChildTexts("/item/pubDate")

		for i, _ := range titles {
			post := PostData{
				Title:       titles[i],
				Description: descriptions[i],
				Link:        links[i],
				Thumbnail:   thumbnails[i],
				PubDate:     pubDates[i],
			}
			posts = append(posts, post)
		}
	})

	if err := rs.c.Visit(rssUrl); nil != err {
		return nil, err
	}

	return posts, nil
}

func NewRssScrapper() *RssScrapper {
	return &RssScrapper{
		c: colly.NewCollector(),
	}
}

func ScrapRssPosts(rssUrl string) []PostData {
	rssScrapper := NewRssScrapper()
	posts, err := rssScrapper.VisitRss(rssUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	return posts
}
