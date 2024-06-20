package vo

import "html/template"

type PostDetails struct {
	Id        int64
	Title     string
	MainImage string
	Content   string
	PubDate   string
}

func (p PostDetails) UnescapedContent() template.HTML {
	return template.HTML(p.Content)
}

func NewPostDetails(id int64, title, mainImage, content, pubDate string) PostDetails {
	return PostDetails{
		Id:        id,
		Title:     title,
		MainImage: mainImage,
		Content:   content,
		PubDate:   pubDate,
	}
}
