package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chepheus/go-rss-parser/web-parser/messanger"
)

const amqpConnStr = "amqp://guest:guest@rabbitmq:5672/"
const queuqName = "rss_post"

func main() {
	amqpMessanger := messanger.NewAMQPMessanger(amqpConnStr, queuqName)
	shutdown := make(chan bool, 1)
	consume(shutdown, amqpMessanger)
	<-shutdown
}

func consume(shutdown chan bool, amqpMessanger messanger.AMQPMessanger) {
	messageCh := amqpMessanger.Consume()
	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case message, ok := <-messageCh:
			if !ok {
				sigKill <- syscall.SIGTERM
				return
			}
			fmt.Println("Message: ", message)
		case <-sigKill:
			amqpMessanger.Close()
			fmt.Println("SIGTERM")
			shutdown <- true
			return
		}
	}
}
