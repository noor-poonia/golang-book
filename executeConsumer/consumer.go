package main

import (
	"fmt"

	"go-rabbitmq/database/rabbitmq"
)

func main()  {
	fmt.Println("RabbitMQ")
	
	rabbitmq.Consumer()
}