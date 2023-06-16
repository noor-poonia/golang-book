package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"encoding/json"
	"go-rabbitmq/model"
	"go-rabbitmq/database/mongodb"

	"github.com/joho/godotenv"
	"github.com/matryer/resync"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct{}

type RabbitConnect struct {
	*amqp.Channel
}

var (
	onceRabbit     resync.Once
	rabbitConn     *amqp.Connection
	instanceRabbit *RabbitConnect
)

func GetRMQConn() *amqp.Connection {
	onceRabbit.Do(func() {
		RMQConnection()
	})

	return rabbitConn
}

func RMQConnection() *amqp.Connection {
	conn, err := amqp.Dial(GetRMQConnectionString())
	if err != nil {
		e := fmt.Sprintf("failed to connect: %s", err)
		panic(e)
	}

	rabbitConn = conn
	return rabbitConn
}

func GetRMQConnectionString() string {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }
	var connect string
	host := os.Getenv("RABBIT_HOST")
	username := os.Getenv("RABBIT_USERNAME")
	password := os.Getenv("RABBIT_PASSWORD")
	port := os.Getenv("RABBIT_PORT")

	connect = fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port)

	return connect
}

func RMQChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		e := fmt.Sprintf("failed to open a channel: %s", err)
		fmt.Println(e)
	}

	return ch
}

func RMQQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(os.Getenv("QUEUE_NAME"), false, false, false, false, nil)
	if err != nil {
		e := fmt.Sprintf("failed to declare a queue: %s", err)
		fmt.Println(e)
	}

	return q
}

func CheckError(err error, msg string) error {
	if err != nil {
		return fmt.Errorf("%s: %s", msg, err)
	}

	return nil
}

// sender
func PublishMessage(book model.Book) error {

	conn := RMQConnection()
	defer conn.Close()
	ch := RMQChannel(conn)
	defer ch.Close()

	queue := RMQQueue(ch)
	
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

// receiver 
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