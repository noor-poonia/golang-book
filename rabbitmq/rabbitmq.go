package rabbitmq

import (
	"fmt"
	"os"

	"github.com/matryer/resync"
	amqp "github.com/rabbitmq/amqp091-go"
	// "go.uber.org/zap"
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

// func (r *RabbitConnect) GetRabbitChannel() *RabbitConnect {
// 	channel, err1 := GetRMQConn().Channel()
// 	if err1 != nil {
// 		zap.L().Fatal("Error in creating channel", zap.Any("error", err1))
// 	}

// 	zap.L().Info("Rabbit MQ connected!")

// 	return &RabbitConnect{channel}
// }

func RMQConnection() *amqp.Connection {
	conn, err := amqp.Dial(getRMQConnectionString())
	if err != nil {
		e := fmt.Sprintf("failed to connect: %s", err)
		fmt.Println(e)
	}

	rabbitConn = conn
	return rabbitConn
}

func getRMQConnectionString() string {
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