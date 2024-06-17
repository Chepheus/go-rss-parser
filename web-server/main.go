package main

import (
	"github.com/Chepheus/go-rss-parser/web-server/pages"
	"github.com/Chepheus/go-rss-parser/web-server/storage"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const port = ":8080"
const connectionStr = "postgresql://root:password@postgres:5432/rss_parser?sslmode=disable"
const migrationSrc = "file://migrations"

func main() {
	db := storage.NewDbConnection(connectionStr)
	storage.MigrationsUp(db, migrationSrc)

	pages.StartWebServer(port, db)
}
