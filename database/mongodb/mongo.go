package mongodb

import (
	"context"
	"fmt"
	"go-rabbitmq/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()
var uri = "mongodb://localhost:27017"
var bookCollection = Collection()

func mongoConn() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err!= nil {
		return nil, err
	}

	return client, nil
}

func Collection() *mongo.Collection {
	if collection == nil {
        client, err := mongoConn()
        if err!= nil {
            panic(err)
        }
    	collection = client.Database("bookdb").Collection("books")
    }
    return collection
}

func InsertBooks(books []model.Book)  {
	for _, book := range books {
		_, err := bookCollection.InsertOne(context.TODO(), book)
		if err!= nil {
			log.Println("failed to insert book")
            panic(err)
        }
		log.Printf("%s inserted.", book.Name)
	}	
}

func GetBooks(query string) ([]model.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var books []model.Book
	if query == "" {
		// TODO: return all books
		cursor, err := bookCollection.Find(context.TODO(), bson.M{})
		defer cursor.Close(ctx)
		if err != nil {
			e := "error while fetching documents"
			return nil, fmt.Errorf(e)
		} else {
			if err := cursor.All(context.Background(), &books); err != nil {
				log.Println("failed to decode books")
				return nil, err
			}
		}
	} else {
		if len(query) != 13 {
			e := "Invalid ISBN"
			return nil, fmt.Errorf(e)
		}
		// TODO: return specifc book
		filter := bson.M{"isbn": query}
		cursor, err := bookCollection.Find(context.TODO(), filter)
		defer cursor.Close(ctx)
		if err != nil {
			e := "failed to find book"
			return nil, fmt.Errorf(e)
		} else {
			if err := cursor.All(context.Background(), &books); err != nil {
				log.Printf("book is: %v\n", books)
				log.Println("Failed to decode books")
				return nil, err
			}
		}
	}
	if len(books) > 0 {
		return books, nil
	}
	e := "No books found"
	return books, fmt.Errorf(e)
}

func UpdateOneBook(query string, book model.Book, updateBook model.Book) (model.Book, error) {  
	var updatedBook model.Book

	update := bson.M{"$set": updateBook}
	filter := bson.M{"isbn":query}

	log.Printf("update value is: %v\n", updateBook)

	options := options.FindOneAndUpdate().SetReturnDocument((options.After))

	err := bookCollection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&updatedBook)
	if err!= nil {
		log.Printf("could not update book: %v\n", err)
		return model.Book{}, err
	}
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
