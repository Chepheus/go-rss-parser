package dto

type WebPost struct {
	ExternalPostLink string
	MainImageSrc     string
	Content          string
}

func NewWebPost(externalPostLink, mainImageSrc, content string) WebPost {
	return WebPost{
		ExternalPostLink: externalPostLink,
		MainImageSrc:     mainImageSrc,
		Content:          content,
	}
}
