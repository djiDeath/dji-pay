package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	. "github.com/tbxark/g4vercel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string
	Email string
}

func Handler(w http.ResponseWriter, r *http.Request) {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://djideath:he9SJ8TGIxKLz4qG@cluster0.y8ltv.mongodb.net/test")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Get a handle for the "users" collection
	collection := client.Database("mydb").Collection("users")

	server := New()

	server.GET("/", func(context *Context) {
		context.JSON(200, H{
			"message": "hello go from vercel !!!!",
		})
	})

	server.GET("/users", func(context *Context) {
		getUsers(context.Writer, context.Req, collection, context)
	})
	server.POST("/users", func(context *Context) {
		createUser(context.Writer, context.Req, collection)
	})

	server.GET("/users", func(context *Context) {

	})

	server.GET("/hello", func(context *Context) {
		name := context.Query("name")
		if name == "" {
			context.JSON(400, H{
				"message": "name not found",
			})
		} else {
			context.JSON(200, H{
				"data": fmt.Sprintf("Hello %s!", name),
			})
		}
	})
	server.GET("/user/:id", func(context *Context) {
		context.JSON(400, H{
			"data": H{
				"id": context.Param("id"),
			},
		})
	})
	server.GET("/long/long/long/path/*test", func(context *Context) {
		context.JSON(200, H{
			"data": H{
				"url": context.Path,
			},
		})
	})
	server.Handle(w, r)
}

func getUsers(w http.ResponseWriter, r *http.Request, collection *mongo.Collection, c *Context) {
	// Find all users
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Convert cursor to []User
	var users []User
	for cursor.Next(context.Background()) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			http.Error(w, "Failed to decode user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Convert []User to JSON and write to response
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	c.JSON(200, H{
		"data": H{
			"url": jsonBytes,
		},
	})
}

func createUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Parse request body into User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Insert user into database
	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	// Return ID of inserted user
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Inserted user with ID: %v", res.InsertedID)))
}
