package main

import (
	"fmt"
	"go-rabbitmq/router"
	"log"
	"net/http"
	"os"

	// "go.uber.org/zap"
	"github.com/joho/godotenv"
)

func main()  {
	fmt.Println("RabbitMQ")

	fmt.Println("loading environment variables")
	err := godotenv.Load()
	port := os.Getenv("MYSQL_PORT")
	fmt.Println("loading finished")
	fmt.Println("PORT: ", port)
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }

	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}

// func init()  {
// 	initLogger()
// }

// func initLogger() {
// 	var logInstance *zap.Logger
// 	var err error
// 	logInstance, err = zap.NewProduction()
// 	if err != nil {
// 		log.Println("Could not create logger instance")
// 		panic(err)
// 	}
// 	if os.Getenv("LOG_LEVEL") == "debug" {
// 		logInstance, err = zap.NewDevelopment()
// 		if err != nil {
// 			log.Println("Could not create logger instance")
// 			panic(err)
// 		}
// 	}
// 	zap.ReplaceGlobals(logInstance)
// }