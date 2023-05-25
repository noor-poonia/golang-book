package controller

import (
	"encoding/json"
	"fmt"
	"go-rabbitmq/model"
	"go-rabbitmq/rabbitmq"
	"go-rabbitmq/validate"
	"net/http"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var book model.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return 
	}

	e := validate.Validate(book)
	if e != nil {
		fmt.Printf("error is: %v\n", e)
		http.Error(w, "Data validation failed", http.StatusBadRequest)
		return
	}

	conn := rabbitmq.RMQConnection()
	defer conn.Close()
	ch := rabbitmq.RMQChannel(conn)
	defer ch.Close()

	queue := rabbitmq.RMQQueue(ch)

	err = rabbitmq.PublishMessage(ch, queue, book)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return 
	}
	fmt.Println("Msg was", book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}