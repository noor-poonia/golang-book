package rabbitmq

import (
	"encoding/json"
	"go-rabbitmq/model"
	"go-rabbitmq/mongodb"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeMessage(ch *amqp.Channel, queue amqp.Queue) (<-chan amqp.Delivery, error) {
	data, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	CheckError(err, "Failed to consume message")

	return data, nil
	
}

func Consumer() {
	conn := RMQConnection()
	defer conn.Close()
	ch := RMQChannel(conn)
	defer ch.Close()

	queue := RMQQueue(ch)

	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// CheckError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// CheckError(err, "Failed to open a channel")
	// defer ch.Close()

	// queue, err := ch.QueueDeclare(
	// 	"books", // name
	// 	false,   // durable
	// 	false,   // delete when unused
	// 	false,   // exclusive
	// 	false,   // no-wait
	// 	nil,     // arguments
	// )
	// CheckError(err, "Failed to declare a queue")
	  
	bookData, err := ConsumeMessage(ch, queue)
	if err!= nil {
        CheckError(err, "failed to consume")
    }
	var books []model.Book
	// go func ()  {
		for d := range bookData {
			// fmt.Printf("d is: %v\n", d)
			var book model.Book
			err := json.Unmarshal(d.Body, &book)
			if err != nil {
				log.Printf("Failed to unmarshal messages: %v\n", err)
				continue
			}
			books = append(books, book)
			mongodb.InsertBooks(books)
			// d.Ack(false)
		}
	// }()
}