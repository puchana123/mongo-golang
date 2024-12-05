package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"mongo-golang/models"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	// validate id as object
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusNotFound)
		return
	}
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// query user form db's collection
	collection := uc.client.Database("mongo-golang").Collection("users")
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// convert data to json
	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
		return
	}

	// respond with user JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// use models user
	var user models.User
	// decode the request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// Assign a new objectID
	user.Id = primitive.NewObjectID()
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Insert into the database
	collection := uc.client.Database("mongo-golang").Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}
	// convert to JSON respond
	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error encoding user data", http.StatusInternalServerError)
		return
	}
	// send message back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	// Validate id as ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusNotFound)
		return
	}
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// remove id from db's collection
	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil || result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// responsd deleted
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s\n", id)
}
