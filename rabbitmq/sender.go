package rabbitmq

import (
	"encoding/json"
	"go-rabbitmq/model"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage(ch *amqp.Channel, queue amqp.Queue, book model.Book) error {
	
	bookBytes, err := json.Marshal(book) 

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: bookBytes,
	})
	
	CheckError(err, "Failed to publish message")
	log.Println("Message Published Successfully")
	log.Printf("message was: %v\n", string(bookBytes))
	return nil
}