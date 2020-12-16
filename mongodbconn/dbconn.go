package mongodbconn

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB(ctx context.Context) *mongo.Client {
	uri := "mongodb://localhost:27017"
	// connect to MongoDB
	mt, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := mt.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected mongodb and pinged.")
	return mt
}

func DisconnectMongodb(ctx context.Context, mt *mongo.Client) {
	log.Println("disconnecting from mongodb")
	fmt.Println("Closing MongoDB Connection...")
	if err := mt.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
