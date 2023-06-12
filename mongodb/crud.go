package mongodb

import (
	"context"
	"go-rabbitmq/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

var bookCollection = Collection()

func InsertBooks(books []model.Book)  {
	for _, book := range books {
		_, err := bookCollection.InsertOne(context.TODO(), book)
		if err!= nil {
            panic(err)
			log.Println("failed to insert book")
        }
		log.Printf("%s inserted.", book.Name)
	}	
}

/// get books => 
func GetBooks(query string)  {
	
}

func GetAllBooks() ([]model.Book, error) {
	cursor, err := bookCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("failed to retrieve books")
		return nil, err
	}
	defer cursor.Close(context.Background())

	var books []model.Book
	if err := cursor.All(context.Background(), &books); err != nil {
		log.Println("failed to decode books")
		return nil, err
	}
	log.Println("retrieved books")

	return books, nil
    
}

func GetOneBook(query string) (model.Book, error)  {
	var book model.Book
	err := bookCollection.FindOne(context.Background(), bson.M{"isbn": query}).Decode(&book)
	if err!= nil {
		log.Println("failed to retrieve book")
		return book, err
	}
	log.Printf("retrieved book: %v\n", book.Name)
	return book, nil
}

func UpdateOneBook(query string, book model.Book) (model.Book, error) {  
	var updatedBook model.Book
	_, e := GetOneBook(query)
	if e != nil {
		log.Printf("Failed to update the book with ISBN: %v\n", query)
	}
	
	update := bson.M{"$set": book}

	err := bookCollection.FindOneAndUpdate(context.Background(), bson.M{"isbn":query}, update).Decode(&updatedBook)
	log.Printf("%v: book updated\n", book.Name)
	return updatedBook, err

}

func DeleteOneBook(isbn string) (model.Book, error) {
	var book model.Book
	filter := bson.M{"isbn": isbn}		
	err := bookCollection.FindOneAndDelete(context.Background(), filter).Decode(&book)
	if err != nil {
		log.Printf("Failed to delete the book: %v\n", err)
		return book, err
	}
	log.Printf("%v: book deleted\n", book.Name)
	return book, nil
}