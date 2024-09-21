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

// getUser function from the previous example
func getUser(db *mongo.Database, email *string, id *string, includePassword bool) (bson.M, error) {
	collection := db.Collection("users")
	var user bson.M
	var filter bson.M

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

	if !includePassword {
		delete(user, "password")
	}

	return user, nil
}

// getArtists fetches all artists and their corresponding user info
func getArtists(db *mongo.Database) ([]bson.M, error) {
	collection := db.Collection("artists")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var artists []bson.M
	for cursor.Next(context.TODO()) {
		var artist bson.M
		if err := cursor.Decode(&artist); err != nil {
			return nil, err
		}

		// Fetch additional user information based on the artist's _id
		artistID := artist["_id"].(primitive.ObjectID).Hex()
		userInfo, err := getUser(db, nil, &artistID, false)
		if err != nil {
			log.Printf("Error fetching user info for artist ID %s: %v", artistID, err)
			continue
		}

		// Merge artist and user info
		for k, v := range userInfo {
			artist[k] = v
		}
		artists = append(artists, artist)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return artists, nil
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

	artists, err := getArtists(db)
	if err != nil {
		log.Fatalf("Error fetching artists: %v", err)
	}

	for _, artist := range artists {
		fmt.Println(artist)
	}
}
