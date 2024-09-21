package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func getUser(db *mongo.Database, email *string, id *string, includePassword bool) (bson.M, error) {
	collection := db.Collection("users")
	var user bson.M
	var filter bson.M

	// Build filter based on email or _id
	if email != nil {
		filter = bson.M{"email": *email}
	} else if id != nil {
		objID, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return nil, err
		}
		filter = bson.M{"_id": objID}
	} else {
		return nil, fmt.Errorf("either email or _id must be provided")
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Optionally remove password field
	if !includePassword {
		delete(user, "password")
	}

	return user, nil
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

	email := "user@example.com"
	user, err := getUser(db, &email, nil, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}
