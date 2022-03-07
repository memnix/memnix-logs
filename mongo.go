package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func Connect() error {

	// Set client options
	clientOptions := options.Client().ApplyURI(MongoURL)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Connect to MongoDB
	client, e := mongo.Connect(ctx, clientOptions)
	if e != nil {
		return e
	}

	// Check the connection
	e = client.Ping(ctx, nil)
	if e != nil {
		return e
	}

	fmt.Println("Connected to mongoDB !")
	// get collection as ref
	db := client.Database("memnix")

	mg = MongoInstance{Client: client, Db: db}

	return nil
}
