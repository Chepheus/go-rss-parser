package vo

type PostListElem struct {
	Id               int64
	Title            string
	Thumbnail        string
	ShortDescription string
	PubDate          string
}

func NewPostListElem(id int64, title, thumbnail, shortDescription, pubDate string) PostListElem {
	return PostListElem{
		Id:               id,
		Title:            title,
		Thumbnail:        thumbnail,
		ShortDescription: shortDescription,
		PubDate:          pubDate,
	}
}
