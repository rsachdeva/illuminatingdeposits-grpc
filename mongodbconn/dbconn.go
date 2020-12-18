package mongodbconn

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongoDB(ctx context.Context, timeoutSec int) *mongo.Client {
	uri := "mongodb://127.0.0.1:27017"
	if host, ok := os.LookupEnv("DEPOSITS_GRPC_DB_HOST"); ok {
		uri = fmt.Sprintf("mongodb://%v:27017", host)
	}
	log.Printf("mongodb uri is %v\n", uri)
	// connect to MongoDB
	// ClientOptions.SetServerSelectionTimeout
	mt, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(time.Duration(timeoutSec)*time.Second))
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
