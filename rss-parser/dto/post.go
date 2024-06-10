package dto

type Post struct {
	Id               int
	Title            string
	ShortDescription string
	ExternalPostLink string
	Thumbnail        string
	PubDate          string
}

func NewPost(title, desc, link, thumbnail, date string) *Post {
	return &Post{
		Title:            title,
		ShortDescription: desc,
		ExternalPostLink: link,
		Thumbnail:        thumbnail,
		PubDate:          date,
	}
}
