package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/Chepheus/go-rss-parser/rss-parser/dto"
	"github.com/Chepheus/go-rss-parser/rss-parser/messanger"
	"github.com/Chepheus/go-rss-parser/rss-parser/scrapper"
	"github.com/Chepheus/go-rss-parser/rss-parser/storage"
	"github.com/Chepheus/go-rss-parser/rss-parser/storage/repository"
)

const period = 5
const rssUrl = "https://feeds.bbci.co.uk/news/world/rss.xml"
const connectionStr = "postgresql://root:password@postgres:5432/rss_parser?sslmode=disable"
const migrationSrc = "file://migrations"

const amqpConnStr = "amqp://guest:guest@rabbitmq:5672/"
const queuqName = "rss_post"

func main() {
	db := storage.NewDbConnection(connectionStr)
	storage.MigrationsUp(db, migrationSrc)

	amqpMessanger := messanger.NewAMQPMessanger(amqpConnStr, queuqName)

	shutdown := make(chan bool, 1)
	go runTicker(shutdown, db, amqpMessanger)
	<-shutdown
}

// Run starts service
func runTicker(shutdown chan bool, db *sql.DB, amqpMessanger messanger.AMQPMessanger) {
	ticker := time.NewTicker(period * time.Second)
	postRepository := repository.NewPostRepository(db)
	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				sigKill <- syscall.SIGTERM
				return
			}
			posts := scrapper.ScrapRssPosts(rssUrl)
			for _, p := range posts {
				post := *dto.NewPost(p.Title, p.Description, p.Link, p.Thumbnail, p.PubDate)
				err := postRepository.Save(post)
				if err != nil {
					log.Fatal(err)
				}

				amqpMessanger.Publish(post.ExternalPostLink)
			}
		case <-sigKill:
			ticker.Stop()
			db.Close()
			amqpMessanger.Close()
			fmt.Println("SIGTERM")
			shutdown <- true
			return
		}
	}
}
