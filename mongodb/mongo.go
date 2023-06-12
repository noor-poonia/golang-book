package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()
var uri = "mongodb://localhost:27017"

func MongoConn() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err!= nil {
		return nil, err
	}

	return client, nil
}

func Collection() *mongo.Collection {
	if collection == nil {
        client, err := MongoConn()
        if err!= nil {
            panic(err)
        }
    	collection = client.Database("bookdb").Collection("books")
    }

	// indexModel := mongo.IndexModel(
	// 	Keys: bson.M{"name": 1},
	// 	Options: options.Index().SetUnique(true),
	// )

	// _, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	// if err != nil {
	// 	log.Printf("Failed to create unique index: %v\n", err)
	// }

    return collection
}