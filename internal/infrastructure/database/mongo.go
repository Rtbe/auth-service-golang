package database

import (
	"context"
	"fmt"
	"log"

	"example.com/auth-service-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//NewMongoClient returns a new mongoDB client.
func NewMongoClient(ctx context.Context, cfg *config.Config) (*mongo.Client, context.Context) {

	connStr := fmt.Sprintf("mongodb://%s:%s@localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0&readPreference=primary&ssl=false", cfg.DbUser, cfg.DbPassword)

	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to MongoDB")
	return client, ctx
}
