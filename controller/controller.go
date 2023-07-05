package controller

import (
	"encoding/json"
	"fmt"
	"go-rabbitmq/constants"
	"go-rabbitmq/model"
	"go-rabbitmq/database/mongodb"
	"go-rabbitmq/database/rabbitmq"
	"go-rabbitmq/utils"
	"log"
	"net/http"
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
			Message: constants.Logs["A105"],
			Data: []model.Book{},
			Error: []string{"Failed to parse data - Check the values"},
		}
		json.NewEncoder(w).Encode(res)
		return
	} else {
		msg, value  := validate.Validate(book)
		fmt.Printf("string recieved from validation: %s\n", value)
		log.Printf("msg is %v\n", value)
		if value != nil {
			log.Printf("error is: %v\n", value)
			res = model.Response{
				Code: "A106",
				Message: constants.Logs["A106"],
				Data: []model.Book{},
				Error: []string{msg},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		} else {
			// TODO: check if book already exists
			b, err := mongodb.GetBooks(book.ISBN)
			if b != nil {
				res = model.Response{
					Code: "A115",
					Message: constants.Logs["A115"],
					Data: []model.Book{},
					Error: []string{"Book already exists with the provided ISBN"},
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)
			} else {
				err = rabbitmq.PublishMessage(book)
				if err != nil {
					res = model.Response{
						Code: "A114",
						Message: constants.Logs["A114"],
						Data: []model.Book{},
						Error: []string{"Check RabbitMQ Connection"},
					}
					w.WriteHeader(http.StatusBadGateway)
					json.NewEncoder(w).Encode(res)
				} else {
					w.WriteHeader(http.StatusAccepted)
					res = model.Response{
						Code: "A101",
						Message: constants.Logs["A101"],
						Data: []model.Book{},
						Error: []string{},
					}
					json.NewEncoder(w).Encode(res)
				}
			}
			
		}
	}
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	res := model.Response{}
	isbn := r.URL.Query().Get("isbn")
	log.Printf("isbn: %v\n", isbn)
	books, err := mongodb.GetBooks(isbn)
	if err!= nil {
    	w.WriteHeader(http.StatusInternalServerError)
		res = model.Response{
			Code: "A104",
			Message: constants.Logs["A104"],
			Data: []model.Book{},
			Error: []string{err.Error()},
		}
    	json.NewEncoder(w).Encode(res)
    } else {
		if books == nil {	
			res = model.Response{
				Code: "A102",
				Message: constants.Logs["A102"],
				Data: []model.Book{},
				Error: []string{"No books found with the given ISBN number"},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		} else {
			res = model.Response{
				Code: "A101",
				Message: constants.Logs["A101"],
				Data: books,
				Error: []string{},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}
}

func UpdateOneBook(w http.ResponseWriter, r *http.Request)  {
	var book model.Book
	res := model.Response{}
	err := json.NewDecoder(r.Body).Decode(&book)
	isbn := book.ISBN
	log.Printf("length isbn: %v\n", len(isbn))
	log.Printf("book is: %v\n", book)
	log.Printf("isbn inside update one book controller function: %v\n", isbn)
	log.Printf("err inside update one book controller function: %v\n", err)
	// checking the data type - parsing error
	if err != nil {
		log.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		res = model.Response{
			Code: "A105",
			Message: constants.Logs["A105"],
			Data: []model.Book{},
			Error: []string{"Failed to parse data - Check the values"},
		}
		json.NewEncoder(w).Encode(res)
		return
	}
	// main update functionality 
	if isbn == "" || len(isbn) != 13 {	
		res = model.Response{
			Code: "A108",
			Message: constants.Logs["A108"],
			Data: []model.Book{},
			Error: []string{"ISBN Number Not Provided or ISBN is not of 13 digits"},
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
	} else {
		b, e := mongodb.GetBooks(isbn)
		log.Printf("Error %v\n", e)
		if e != nil {
			// log.Printf("no book found: %v\n", err.Error())
			res = model.Response{
				Code: "A102",
				Message: "No book found with the provided ISBN",
				Data: []model.Book{},
				Error: []string{"Book not found with the provided ISBN, Check your ISBN number again"},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		} else {
			msg, value := validate.Validate(book)
			log.Printf("value: %v\n", value)
			if value != nil {
				fmt.Printf("error is: %v\n", e)
				res = model.Response{
					Code: "A106",
					Message: constants.Logs["A106"],
					Data: []model.Book{},
					Error: []string{msg},
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(res)

			} else {
				updatedBook, er := mongodb.UpdateOneBook(isbn, b[0], book)
				log.Printf("value of updated book is: %v\n", updatedBook)
				if er != nil {
					log.Printf("book not updated: %v\n", er)
					w.WriteHeader(http.StatusBadGateway)
					res = model.Response{
						Code: "A110",
						Message: constants.Logs["A110"],
						Data: []model.Book{},
						Error: []string{"Failed to Update the Book with the given ISBN"},
					}
					json.NewEncoder(w).Encode(res)
				} else {
					w.WriteHeader(http.StatusOK)
					res = model.Response{
						Code: "A109",
						Message: constants.Logs["A109"],
						Data: []model.Book{updatedBook},
						Error: []string{},
					}
					json.NewEncoder(w).Encode(res)
				}
			}
		}
	}
}

func DeleteOneBook(w http.ResponseWriter, r *http.Request)  {
	isbn := r.URL.Query().Get("isbn")
	res := model.Response{}
	if isbn == "" {
		res = model.Response{
			Code: "A108",
			Message: constants.Logs["A108"],
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
				Message: constants.Logs["A102"],
				Data: []model.Book{},
				Error: []string{"No Book Found with the requested ISBN"},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		} else {
			res = model.Response{
				Code: "A107",
				Message: constants.Logs["A107"],
				Data: []model.Book{book},
				Error: []string{},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}
}