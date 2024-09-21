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

// createEvent inserts a new event into the "events" collection
func createEvent(db *mongo.Database, createdBy, location, genre, description string, artists []string, eventDate time.Time) (primitive.ObjectID, error) {
	collection := db.Collection("events")

	event := bson.M{
		"created_by": createdBy,
		"location":   location,
		"description": description,
		"genre":      genre,
		"artists":    artists,
		"date":       eventDate,
		"created":    time.Now().UTC(),
	}

	result, err := collection.InsertOne(context.TODO(), event)
	if err != nil {
		return primitive.NilObjectID, err
	}

	eventID := result.InsertedID.(primitive.ObjectID)
	fmt.Println("Event ID:", eventID)
	return eventID, nil
}

// getEvents fetches all events from the "events" collection
func getEvents(db *mongo.Database) ([]bson.M, error) {
	collection := db.Collection("events")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var events []bson.M
	if err = cursor.All(context.TODO(), &events); err != nil {
		return nil, err
	}

	for _, event := range events {
		fmt.Println(event)
	}

	return events, nil
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

	// Example: Create an event
	artists := []string{"Artist 1", "Artist 2"}
	eventDate := time.Now().Add(48 * time.Hour)
	eventID, err := createEvent(db, "user@example.com", "New York", "Jazz", "A cool jazz event", artists, eventDate)
	if err != nil {
		log.Fatalf("Error creating event: %v", err)
	}
	fmt.Println("Created event with ID:", eventID)

	// Example: Get all events
	events, err := getEvents(db)
	if err != nil {
		log.Fatalf("Error fetching events: %v", err)
	}
	fmt.Println("Fetched events:", events)
}
