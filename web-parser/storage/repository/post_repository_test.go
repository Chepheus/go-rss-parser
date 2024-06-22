package repository

import (
	"database/sql"
	"log"
	"testing"

	"github.com/Chepheus/go-rss-parser/web-parser/dto"
	"github.com/DATA-DOG/go-sqlmock"
)

var post = dto.WebPost{
	MainImageSrc:     "https://test.com/image.jpg",
	Content:          "test",
	ExternalPostLink: "https://test.com",
}

func TestUpdate(t *testing.T) {
	db, mock := NewMock()

	repo := NewPostRepository(db)

	mock.
		ExpectExec("UPDATE posts").
		WithArgs(post.MainImageSrc, post.Content, post.ExternalPostLink).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(post)
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
