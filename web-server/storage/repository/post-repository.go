package repository

import (
	"database/sql"
	"fmt"
)

type Post struct {
	Id               int64
	Title            string
	Thumbnail        sql.NullString
	MainImage        sql.NullString
	ShortDescription sql.NullString
	Content          sql.NullString
	PubDate          string
}

type PostRepository struct {
	db *sql.DB
}

func (r *PostRepository) GetPosts(page, limit int32) ([]Post, error) {
	rows, err := r.db.Query(
		"SELECT id, title, thumbnail, short_description, pub_date FROM posts ORDER BY pub_date LIMIT $1 OFFSET $2",
		limit,
		(page-1)*limit,
	)
	defer func() {
		_ = rows.Close()
	}()

	if err != nil {
		return nil, fmt.Errorf("page %d is't exists", page)
	}

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Id, &post.Title, &post.Thumbnail, &post.ShortDescription, &post.PubDate); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPost(postId int64) (*Post, error) {
	var post Post
	err := r.db.QueryRow(
		"SELECT id, title, main_image, content, pub_date FROM posts WHERE id = $1",
		postId,
	).Scan(&post.Id, &post.Title, &post.MainImage, &post.Content, &post.PubDate)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}
