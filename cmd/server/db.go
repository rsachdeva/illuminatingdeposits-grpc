package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectMongoDB() (context.Context, *mongo.Client) {
	uri := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	return ctx, mt
}

func disconnectMongodb(mt *mongo.Client, ctx context.Context) {
	log.Println("disconnecting from mongodb")
	fmt.Println("Closing MongoDB Connection...")
	if err := mt.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
