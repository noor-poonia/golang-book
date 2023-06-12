package router

import (
	"go-rabbitmq/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/books", controller.GetBooks).Methods("GET")
	r.HandleFunc("/books", controller.DeleteOneBook).Methods("DELETE")
	r.HandleFunc("/books", controller.UpdateOneBook).Methods("PUT")
	r.HandleFunc("/book", controller.CreateMessage).Methods("POST")
	
	return r
}