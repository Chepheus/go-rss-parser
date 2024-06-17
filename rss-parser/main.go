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

const period = 10
const rssUrl = "https://feeds.bbci.co.uk/news/world/rss.xml"
const connectionStr = "postgresql://root:password@postgres:5432/rss_parser?sslmode=disable"

const amqpConnStr = "amqp://guest:guest@rabbitmq:5672/"
const queuqName = "rss_post"

func main() {
	db := storage.NewDbConnection(connectionStr)

	amqpMessanger := messanger.NewAMQPMessanger(amqpConnStr, queuqName)
	postRepository := repository.NewPostRepository(db)

	shutdown := make(chan bool, 1)
	go runTicker(shutdown, db, amqpMessanger, postRepository)
	<-shutdown
}

// Run starts service
func runTicker(shutdown chan bool, db *sql.DB, amqpMessanger messanger.AMQPMessanger, postRepository repository.PostRepository) {
	ticker := time.NewTicker(period * time.Second)
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
				go func(p scrapper.PostData) {
					isExist, err := postRepository.IsExist(p.Link)
					if err != nil {
						log.Fatal(err)
					}

					if !isExist {
						post := *dto.NewPost(p.Title, p.Description, p.Link, p.Thumbnail, p.PubDate)
						err := postRepository.Save(post)
						if err != nil {
							log.Fatal(err)
						}

						err = amqpMessanger.Publish(post.ExternalPostLink)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println("[rss-parser] Published to RebbitMQ: " + post.ExternalPostLink)
					}
				}(p)
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
