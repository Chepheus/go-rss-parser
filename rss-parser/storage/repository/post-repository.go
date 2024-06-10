package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chepheus/go-rss-parser/rss-parser/dto"
)

type PostRepository struct {
	db *sql.DB
}

func (r *PostRepository) Save(post dto.Post) error {
	fmt.Println(post)
	_, err := r.db.Exec(
		"INSERT INTO posts (title, short_description, external_post_link, thumbnail, pub_date) VALUES($1, $2, $3, $4, $5) ON CONFLICT (title) DO NOTHING;",
		post.Title,
		post.ShortDescription,
		post.ExternalPostLink,
		post.Thumbnail,
		post.PubDate,
	)

	if err != nil {
		return errors.New("post with title was not saved: " + post.Title)
	}

	return nil
}

func NewPostRepository(db *sql.DB) PostRepository {
	return PostRepository{db}
}
