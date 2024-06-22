package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/Chepheus/go-rss-parser/rss-parser/dto"
	"github.com/DATA-DOG/go-sqlmock"
)

var post = dto.Post{
	Title:            "test",
	ShortDescription: "test",
	ExternalPostLink: "https://test.com",
	Thumbnail:        "htpps://test.com/img.jpg",
	PubDate:          "Thu, 20 Jun 2024 15:57:36 GMT",
}

func TestIsExist(t *testing.T) {
	db, mock := NewMock()

	repo := NewPostRepository(db)

	query := "SELECT COUNT\\(\\*\\) FROM posts WHERE external_post_link = \\$1"

	rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(true)

	mock.ExpectQuery(query).WithArgs(post.ExternalPostLink).WillReturnRows(rows)

	isExist, err := repo.IsExist(post.ExternalPostLink)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !isExist {
		t.Error("row isn't exist")
	}
}

func TestSave(t *testing.T) {
	db, mock := NewMock()

	repo := NewPostRepository(db)

	mock.
		ExpectExec("INSERT INTO posts").
		WithArgs(post.Title, post.ShortDescription, post.ExternalPostLink, post.Thumbnail, post.PubDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save(post)
	if err != nil {
		t.Error(err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
