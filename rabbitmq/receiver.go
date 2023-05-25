package rabbitmq

import (
	"encoding/json"
	"fmt"
	"go-rabbitmq/model"
	"go-rabbitmq/mysql"
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

func Consumer()  {
	conn := RMQConnection()
	defer conn.Close()
	ch := RMQChannel(conn)
	defer ch.Close()

	queue := RMQQueue(ch)

	bookData, err := ConsumeMessage(ch, queue)
	if err!= nil {
        CheckError(err, "failed to consume")
        return
    }
	mysql.MysqlConnction()
	var books []model.Book
	for d := range bookData {
		fmt.Printf("d is: %v\n", string(d.Body))
		var book model.Book
		err := json.Unmarshal(d.Body, &book)
		if err != nil {
			log.Printf("Failed to unmarshal messages: %v\n", err)
			continue
		}
		books = append(books, book)
		fmt.Printf("books: %v\n", books)
		// d.Ack(false)
	}
	log.Println(books)
}