package controller

import (
	"encoding/json"
	"fmt"
	"go-rabbitmq/logs"
	"go-rabbitmq/model"
	"go-rabbitmq/mongodb"
	"log"
	"go-rabbitmq/rabbitmq"
	"go-rabbitmq/validate"
	"net/http"

	// "github.com/gorilla/mux"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	res := model.Response{}
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		res = model.Response{
			Code: "A105",
			Message: logs.Logs["A105"],
			Data: []model.Book{},
			Error: []string{err.Error()},
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	e := validate.Validate(book)
	if e != nil {
		fmt.Printf("error is: %v\n", e)
		res = model.Response{
			Code: "A106",
			Message: logs.Logs["A106"],
			Data: []model.Book{},
			Error: []string{e.Error()},
		}
		json.NewEncoder(w).Encode(res)
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
	w.WriteHeader(http.StatusAccepted)
	res = model.Response{
		Code: "A101",
		Message: logs.Logs["A101"],
		Data: []model.Book{},
		Error: []string{},
	}
	json.NewEncoder(w).Encode(res)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	res := model.Response{}
	// books := []model.Book{}
	isbn := r.URL.Query().Get("isbn")
	log.Printf("isbn: %v\n", isbn)
	if isbn == "" {
		// TODO: return all the books
		books, err := mongodb.GetAllBooks()
		if err!= nil {
    		w.WriteHeader(http.StatusInternalServerError)
			res = model.Response{
				Code: "A104",
				Message: logs.Logs["A104"],
				Data: []model.Book{},
				Error: []string{"Failed to retriev books - check mongodb connection"},
			}
    		json.NewEncoder(w).Encode(res)
    	}
		res = model.Response{
			Code: "A101",
			Message: logs.Logs["A101"],
			Data: books,
			Error: []string{},
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(res)
	} else {
		// TODO: return the book
		// var book model.Book
		book, err := mongodb.GetOneBook(isbn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res = model.Response{
				Code: "A102",
				Message: logs.Logs["A102"],
				Data: []model.Book{},
				Error: []string{err.Error()},
			}
			json.NewEncoder(w).Encode(res)
		} else {
			res = model.Response{
				Code: "A101",
				Message: logs.Logs["A101"],
				Data: []model.Book{book},
				Error: []string{},
			}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(res)
		}
	}
}

func UpdateOneBook(w http.ResponseWriter, r *http.Request)  {
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	isbn := book.ISBN
	// isbn := r.URL.Query().Get("isbn")
	res := model.Response{}
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		res = model.Response{
			Code: "A105",
			Message: logs.Logs["A105"],
			Data: []model.Book{},
			Error: []string{err.Error()},
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	if isbn == "" {	
		res = model.Response{
			Code: "A108",
			Message: logs.Logs["A108"],
			Data: []model.Book{},
			Error: []string{"ISBN Number Not Provided"},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
	} else {
		
		

		e := validate.Validate(book)
		if e != nil {
			fmt.Printf("error is: %v\n", e)
			res = model.Response{
				Code: "A106",
				Message: logs.Logs["A106"],
				Data: []model.Book{},
				Error: []string{e.Error()},
			}
			json.NewEncoder(w).Encode(res)
			return
		}
		b, er := mongodb.UpdateOneBook(isbn, book)
		if er != nil {
			log.Printf("book not updated: %v\n", er)
			w.WriteHeader(http.StatusBadGateway)
			res = model.Response{
				Code: "A110",
				Message: logs.Logs["A110"],
				Data: []model.Book{},
				Error: []string{"Failed to Update the Book with the given ISBN"},
			}
			json.NewEncoder(w).Encode(res)
		} else {
			w.WriteHeader(http.StatusAccepted)
			res = model.Response{
				Code: "A109",
				Message: logs.Logs["A109"],
				Data: []model.Book{b},
				Error: []string{},
			}
			json.NewEncoder(w).Encode(res)
		}
	}
}

func DeleteOneBook(w http.ResponseWriter, r *http.Request)  {
	isbn := r.URL.Query().Get("isbn")
	res := model.Response{}
	if isbn == "" {
		res = model.Response{
			Code: "A108",
			Message: logs.Logs["A108"],
			Data: []model.Book{},
			Error: []string{"No ISBN Found in the Query Parameters"},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
	} else {
		book, err := mongodb.DeleteOneBook(isbn)
		if err != nil {
			res = model.Response{
				Code: "A102",
				Message: logs.Logs["A102"],
				Data: []model.Book{},
				Error: []string{"No Book Found with the requested ISBN"},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		} else {
			res = model.Response{
				Code: "A107",
				Message: logs.Logs["A107"],
				Data: []model.Book{book},
				Error: []string{},
			}
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(res)
		}
	}
}