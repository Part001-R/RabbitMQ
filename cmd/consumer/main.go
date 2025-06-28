package main

import (
	"log"
	"time"

	rmq "github.com/rabbitmq/amqp091-go"
)

// failOnError - вспомогательная функция для обработки ошибок
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// 1. Connecting RabbitMQ
	conn, err := rmq.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// 2. Connecting channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 3. Declare queue
	q, err := ch.QueueDeclare(
		"que_message", // name queue
		false,         // durable: saving queue after restart RabbitMQ
		false,         // delete when unused: if true - delete queue if missed consumers
		false,         // exclusive: Если true, очередь будет доступна только для текущего соединения и будет удалена при его закрытии
		false,         // no-wait: Если true, клиент не будет ждать ответа от сервера
		nil,           // arguments: Дополнительные аргументы для очереди
	)
	failOnError(err, "Failed to declare a queue")

	// 4. Регистрация потребителя для получения сообщений
	msgs, err := ch.Consume(
		q.Name, // queue name:
		"",     // consumer: Уникальный идентификатор потребителя (пустая строка генерирует уникальный)
		false,  // auto-ack:
		// Если auto-ack = false, то сообщения НЕ БУДУТ автоматически подтверждаться RabbitMQ.
		// Потребитель ДОЛЖЕН явно отправить подтверждение (Ack) после обработки.
		false, // exclusive: Если true, только этот потребитель может получать сообщения из этой очереди
		false, // no-local: Если true, сервер не будет отправлять сообщения обратно тому же соединению, которое их опубликовало
		false, // no-wait: Если true, клиент не будет ждать ответа от сервера
		nil,   // args: Дополнительные аргументы для потребителя
	)
	failOnError(err, "Failed to register a consumer")

	// for continue runnig
	forever := make(chan bool)

	// Handler messages
	go func() {
		for d := range msgs {
			log.Printf(" [x] Received: %s", d.Body)
			time.Sleep(2 * time.Second)

			err = d.Ack(false) // false - ACK current message, true - ACK all previe messages
			if err != nil {
				failOnError(err, string(d.Body))
			}

			log.Printf(" [x] Message acknowledged: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
