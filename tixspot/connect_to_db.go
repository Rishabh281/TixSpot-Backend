package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// connectToDB loads the URI from .env, creates a MongoDB client, and connects to the server.
func connectToDB() *mongo.Client {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the MongoDB URI from environment variables
	uri := os.Getenv("URI")
	if uri == "" {
		log.Fatal("URI not found in .env file")
	}

	// Create a new MongoDB client and connect to the server
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Ping the MongoDB deployment to verify connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping the MongoDB server: %v", err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

func main() {
	client := connectToDB()
	// You can now use the client to interact with your MongoDB database
	defer client.Disconnect(context.TODO())
}
