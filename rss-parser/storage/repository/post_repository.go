package repository

import (
	"database/sql"
	"errors"

	"github.com/Chepheus/go-rss-parser/rss-parser/dto"
)

type PostRepository struct {
	db *sql.DB
}

func (r *PostRepository) IsExist(externalPostLink string) (bool, error) {
	var isExist bool
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM posts WHERE external_post_link = $1",
		externalPostLink,
	).Scan(&isExist)

	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (r *PostRepository) Save(post dto.Post) error {
	_, err := r.db.Exec(
		"INSERT INTO posts (title, short_description, external_post_link, thumbnail, pub_date) VALUES($1, $2, $3, $4, $5) ON CONFLICT (external_post_link) DO NOTHING;",
		post.Title,
		post.ShortDescription,
		post.ExternalPostLink,
		post.Thumbnail,
		post.PubDate,
	)

	if err != nil {
		return errors.New("post with link was not saved: " + post.ExternalPostLink)
	}

	return nil
}

func NewPostRepository(db *sql.DB) PostRepository {
	return PostRepository{db}
}
