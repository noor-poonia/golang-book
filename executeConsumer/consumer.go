package main

import (
	"fmt"
	"go-rabbitmq/rabbitmq"
	// "go-rabbitmq/router"
	"log"
	// "net/http"

	"github.com/joho/godotenv"
)

func main()  {
	fmt.Println("RabbitMQ")

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }
	
	// r := router.Router()
	// log.Fatal(http.ListenAndServe(":8080", r))

	rabbitmq.Consumer()
}