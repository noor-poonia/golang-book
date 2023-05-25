package router

import (
	"go-rabbitmq/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	// r.HandleFunc("/books", controller.GetMessages).Methods("GET")
	r.HandleFunc("/book", controller.CreateMessage).Methods("POST")
	
	return r
}