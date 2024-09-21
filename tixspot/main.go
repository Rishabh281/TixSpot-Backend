package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

type CreateEvent struct {
	CreatedBy  string               `json:"created_by" binding:"required"`
	Location   string               `json:"location" binding:"required"`
	Description string              `json:"description" binding:"required"`
	Genre      string               `json:"genre" binding:"required"`
	Artists    []string             `json:"artists"`
	Date       time.Time            `json:"date" binding:"required"`
	Created    time.Time            `json:"created" binding:"required"`
}

var client *mongo.Client

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

func validateToken(token string, tokenType string) (map[string]interface{}, error) {
	// Simulate token validation
	if token == "" {
		return nil, fmt.Errorf("invalid token")
	}
	// Assuming a valid user is retrieved
	user := map[string]interface{}{
		"username": "example_user",
		"email":    "example_user@example.com",
		"password": "hashed_password", // This should be removed when returning
	}
	delete(user, "password") // Ensure password is removed
	return user, nil
}

func createEvent(c *mongo.Collection, event CreateEvent) {
	_, err := c.InsertOne(context.TODO(), event)
	if err != nil {
		log.Fatal(err)
	}
}

func getEvents(c *mongo.Collection) []bson.M {
	var events []bson.M
	cursor, err := c.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var event bson.M
		if err = cursor.Decode(&event); err != nil {
			log.Fatal(err)
		}
		events = append(events, event)
	}
	return events
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	client = connectToDB()

	router.GET("/user/details", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		user, err := validateToken(token, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	router.POST("/events/create", func(c *gin.Context) {
		var createEventForm CreateEvent
		if err := c.ShouldBindJSON(&createEventForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")
		_, err := validateToken(token, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		collection := client.Database("tixspot").Collection("events")
		createEvent(collection, createEventForm)
		c.JSON(http.StatusOK, createEventForm)
	})

	router.GET("/events/getall", func(c *gin.Context) {
		collection := client.Database("tixspot").Collection("events")
		events := getEvents(collection)
		c.JSON(http.StatusOK, events)
	})

	router.GET("/artists", func(c *gin.Context) {
		collection := client.Database("tixspot").Collection("artists")
		var artists []bson.M
		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var artist bson.M
			if err = cursor.Decode(&artist); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			artists = append(artists, artist)
		}

		c.JSON(http.StatusOK, artists)
	})

	router.Run(":8080")
}
