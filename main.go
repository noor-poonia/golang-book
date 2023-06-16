package main

import (
	"fmt"
	"go-rabbitmq/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main()  {
	fmt.Println("RabbitMQ")

	err := godotenv.Load()
    if err != nil {
        log.Fatalf("err loading: %v", err)
    }

	r := mux.NewRouter()

	r.HandleFunc("/books", controller.GetBooks).Methods("GET")
	// r.HandleFunc("/books", controller.GetBook).Methods("GET")
	r.HandleFunc("/books", controller.DeleteOneBook).Methods("DELETE")
	r.HandleFunc("/books", controller.UpdateOneBook).Methods("PUT")
	r.HandleFunc("/book", controller.CreateMessage).Methods("POST")
	
	log.Fatal(http.ListenAndServe(":8080", r))
}