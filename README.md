Learning RabbitMQ.

docker - `docker run -d --hostname my-rabbit --name rabbitmq -p 5672:5672 -p 8080:15672 rabbitmq:3-management`

web - `http://localhost:8080`
login:    `guest`
password: `guest`

Run producer - `go run ./cmd/producer/main.go`
Run consumer - `go run ./cmd/consumer/main.go`

+ `v0.0.1` - Producer, consumer RabbitMQ. 