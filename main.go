package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mongo-golang/controllers"
)

func main() {
	// create MongoDB session
	client, err := getSession()
	if err != nil {
		fmt.Printf("Error initializing MongoDB client: %v\n", err)
		return
	}
	// create new router
	r := httprouter.New()
	// Initialize the UserController
	uc := controllers.NewUserController(client)
	// define routes
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	// Start HTTP server
	fmt.Println("Starting server on localhost: 8080")
	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

func getSession() (*mongo.Client, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connnect to mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MongoDB: %w", err)
	}
	// if get error from connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}
	return client, nil
}
