package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var post = Post{
	Id:    1,
	Title: "test",
	Thumbnail: sql.NullString{
		String: "htpps://test.com/thumbnail.jpg",
		Valid:  false,
	},
	MainImage: sql.NullString{
		String: "htpps://test.com/main_img.jpg",
		Valid:  false,
	},
	ShortDescription: sql.NullString{
		String: "short description",
		Valid:  false,
	},
	Content: sql.NullString{
		String: "post content",
		Valid:  false,
	},
	PubDate: "Thu, 20 Jun 2024 15:57:36 GMT",
}

func TestGetPosts(t *testing.T) {
	db, mock := NewMock()

	repo := NewPostRepository(db)

	query := "SELECT id, title, thumbnail, short_description, pub_date FROM posts ORDER BY pub_date LIMIT \\$1 OFFSET \\$2"

	rows := sqlmock.
		NewRows([]string{"id", "title", "thumbnail", "short_description", "pub_date"}).
		AddRow(post.Id, post.Title, post.Thumbnail.String, post.ShortDescription.String, post.PubDate)

	mock.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)

	posts, err := repo.GetPosts(1, 10)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if len(posts) == 0 {
		t.Error("posts are empty")
	}

	if len(posts) > 0 && posts[0].Id != post.Id {
		t.Errorf("id is't equal %d", post.Id)
	}
}

func TestGetPost(t *testing.T) {
	db, mock := NewMock()

	repo := NewPostRepository(db)

	query := "SELECT id, title, main_image, content, pub_date FROM posts WHERE id = \\$1"

	rows := sqlmock.
		NewRows([]string{"id", "title", "main_image", "content", "pub_date"}).
		AddRow(post.Id, post.Title, post.MainImage.String, post.Content.String, post.PubDate)

	mock.ExpectQuery(query).WithArgs(post.Id).WillReturnRows(rows)

	p, err := repo.GetPost(post.Id)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if p.Id != post.Id {
		t.Errorf("id is't equal %d", post.Id)
	}
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
