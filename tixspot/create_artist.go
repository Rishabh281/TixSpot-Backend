package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createArtist inserts a new artist into the "artists" collection if the artist does not already exist
func createArtist(db *mongo.Database, artistID primitive.ObjectID, stageName, description, genre string) (primitive.ObjectID, error) {
	collection := db.Collection("artists")

	// Check if the artist already exists by ID
	var existingArtist bson.M
	err := collection.FindOne(context.TODO(), bson.M{"_id": artistID}).Decode(&existingArtist)
	if err != mongo.ErrNoDocuments {
		fmt.Println("Artist already exists")
		return primitive.NilObjectID, nil
	}

	// Create the artist document
	artist := bson.M{
		"_id":         artistID,
		"stage_name":  stageName,
		"description": description,
		"genre":       genre,
		"date":        time.Now().UTC(),
	}

	// Insert the artist into the database
	result, err := collection.InsertOne(context.TODO(), artist)
	if err != nil {
		return primitive.NilObjectID, err
	}

	newArtistID := result.InsertedID.(primitive.ObjectID)
	fmt.Println("Artist ID:", newArtistID)
	return newArtistID, nil
}

func connectToDB() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func main() {
	client := connectToDB()
	db := client.Database("your_db_name")

	// Example: Create a new artist
	artistID := primitive.NewObjectID()
	stageName := "Cool Artist"
	description := "A great performer"
	genre := "Jazz"

	newArtistID, err := createArtist(db, artistID, stageName, description, genre)
	if err != nil {
		log.Fatalf("Error creating artist: %v", err)
	}
	fmt.Println("Created artist with ID:", newArtistID)
}
