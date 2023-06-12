package main

import (
	"fmt"

	"go-rabbitmq/rabbitmq"
)

func main()  {
	fmt.Println("RabbitMQ")
	
	rabbitmq.Consumer()
}