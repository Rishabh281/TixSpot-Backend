package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createUser inserts a new user into the "users" collection if the email does not already exist
func createUser(db *mongo.Database, email, password, firstName, lastName, username string) (primitive.ObjectID, error) {
	collection := db.Collection("users")

	// Check if the user already exists by email
	var existingUser bson.M
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		fmt.Println("User already exists")
		return primitive.NilObjectID, nil
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create the user document
	user := bson.M{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"username":   username,
		"password":   string(hashedPassword),
		"date":       time.Now().UTC(),
	}

	// Insert the user into the database
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	userID := result.InsertedID.(primitive.ObjectID)
	fmt.Println("User ID:", userID)
	return userID, nil
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

	// Example: Create a new user
	email := "user@example.com"
	password := "securepassword"
	firstName := "John"
	lastName := "Doe"
	username := "johndoe123"

	userID, err := createUser(db, email, password, firstName, lastName, username)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	fmt.Println("Created user with ID:", userID)
}
