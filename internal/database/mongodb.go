package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMongoCollection(dsn string) (*mongo.Database, error) {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB")
	database := client.Database(os.Getenv("MONGO_DATABASE"))
	if err != nil {
		return nil, err
	}
	return database, nil
}
