package messanger

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPMessanger struct {
	amqpConn *amqp.Connection
	amqpCh   *amqp.Channel
	queue    *amqp.Queue
}

func (m AMQPMessanger) Consume() chan string {
	msgs, err := m.amqpCh.Consume(
		m.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatal(err)
	}

	messageCh := make(chan string)
	go func() {
		for m := range msgs {
			log.Printf("Received a message: %s", m.Body)
			messageCh <- string(m.Body)
		}
	}()
	return messageCh
}

func (m AMQPMessanger) Close() {
	m.amqpCh.Close()
	m.amqpConn.Close()
}

func NewAMQPMessanger(connStr, queuqName string) AMQPMessanger {
	conn := reconnectConnection(connStr)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		queuqName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	return AMQPMessanger{
		amqpConn: conn,
		amqpCh:   ch,
		queue:    &q,
	}
}

func reconnectConnection(connStr string) *amqp.Connection {
	fmt.Println("[amqp]: connection")
	conn, err := amqp.Dial(connStr)
	if err != nil {
		time.Sleep(1 * time.Second)
		conn = reconnectConnection(connStr)
	}

	return conn
}
