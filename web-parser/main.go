package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/Chepheus/go-rss-parser/web-parser/dto"
	"github.com/Chepheus/go-rss-parser/web-parser/messanger"
	"github.com/Chepheus/go-rss-parser/web-parser/scrapper"
	"github.com/Chepheus/go-rss-parser/web-parser/storage"
	"github.com/Chepheus/go-rss-parser/web-parser/storage/repository"
)

const connectionStr = "postgresql://root:password@postgres:5432/rss_parser?sslmode=disable"
const amqpConnStr = "amqp://guest:guest@rabbitmq:5672/"
const queuqName = "rss_post"

func main() {
	db := storage.NewDbConnection(connectionStr)
	amqpMessanger := messanger.NewAMQPMessanger(amqpConnStr, queuqName)

	postRepository := repository.NewPostRepository(db)

	shutdown := make(chan bool, 1)
	go consume(shutdown, db, amqpMessanger, postRepository)
	<-shutdown
}

func consume(shutdown chan bool, db *sql.DB, amqpMessanger messanger.AMQPMessanger, postRepository repository.PostRepository) {
	messageCh := amqpMessanger.Consume()
	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case postLink, ok := <-messageCh:
			if !ok {
				sigKill <- syscall.SIGTERM
				return
			}

			go func() {
				webScrapper := scrapper.NewWebScrapper()
				webPost, err := webScrapper.Visit(postLink)
				if err != nil {
					log.Fatal(err)
				}

				var content string
				for _, v := range webPost.TextBlocks {
					content = content + "<p>" + v + "</p>"
				}

				err = postRepository.Update(dto.NewWebPost(
					webPost.ExternalPostLink,
					webPost.MainImageSrc,
					content,
				))

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("[web-parser] Post content saved to DB: " + postLink)
			}()
		case <-sigKill:
			amqpMessanger.Close()
			db.Close()
			fmt.Println("SIGTERM")
			shutdown <- true
			return
		}
	}
}
