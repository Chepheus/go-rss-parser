package repository

import (
	"database/sql"

	"github.com/Chepheus/go-rss-parser/web-parser/dto"
)

type PostRepository struct {
	db *sql.DB
}

func (r PostRepository) Update(post dto.WebPost) error {
	_, err := r.db.Exec(
		"UPDATE posts SET main_image = $1, content = $2 WHERE external_post_link = $3",
		post.MainImageSrc,
		post.Content,
		post.ExternalPostLink,
	)

	if err != nil {
		return err
	}

	return nil
}

func NewPostRepository(db *sql.DB) PostRepository {
	return PostRepository{
		db: db,
	}
}
