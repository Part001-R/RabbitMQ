package main

import (
	"context"
	"fmt"
	"log"
	"time"

	rmq "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// 1. Connection RabbitMQ
	conn, err := rmq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. Conction channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. Init Queue
	q, err := ch.QueueDeclare(
		"que_message", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 4. Sent message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := fmt.Sprintf("test message {%s}", time.Now().Format("01-02-2006 15-04-05"))
	err = ch.PublishWithContext(ctx,
		"",     // exchange - пустая строка означает 'default exchange'
		q.Name, // routing key - имя очереди, в которую выполняется отправка
		false,  // mandatory
		false,  // immediate
		rmq.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", body)
}
