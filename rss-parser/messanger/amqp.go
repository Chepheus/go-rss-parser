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

func (m AMQPMessanger) Publish(body string) error {
	err := m.amqpCh.Publish(
		"",           // exchange
		m.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	return err
}

func (m AMQPMessanger) Close() {
	m.amqpCh.Close()
	m.amqpConn.Close()
}

func NewAMQPMessanger(connStr, queuqName string) AMQPMessanger {
	amqpConn := reconnectConnection(connStr)

	amqpCh, err := amqpConn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := amqpCh.QueueDeclare(
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
		amqpConn: amqpConn,
		amqpCh:   amqpCh,
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
