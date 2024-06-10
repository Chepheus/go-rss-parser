package storage

import (
	"database/sql"
	"fmt"
	"time"
)

func NewDbConnection(connStr string) *sql.DB {
	db := reconnectDbConnection(connStr)

	return db
}

func reconnectDbConnection(connStr string) *sql.DB {
	fmt.Println("connection")
	db, err := sql.Open("postgres", connStr)
	if pingErr := db.Ping(); err != nil || pingErr != nil {
		fmt.Println(err, pingErr)
		time.Sleep(1 * time.Second)
		db = reconnectDbConnection(connStr)
	}

	return db
}
