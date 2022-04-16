package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yash04g/Go-mux-mongo/configs"
	"github.com/yash04g/Go-mux-mongo/models"
	"github.com/yash04g/Go-mux-mongo/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		// Validating the user body
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Using the validator library validating the required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		newUser := models.User{
			ID:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}
		res, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusCreated)
		response := responses.UserResponse{Status: http.StatusCreated, Message: "New user created", Data: map[string]interface{}{"data": res}}
		json.NewEncoder(w).Encode(response)
	}
}

func GetAUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		params := mux.Vars(r)
		userId := params["userId"]
		var user models.User
		defer cancel()

		objectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fmt.Println("Error while getting object id", err)
			return
		}
		err = userCollection.FindOne(ctx, bson.M{"id": objectId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error in finding the user", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(w).Encode(response)
	}
}

func EditAUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		params := mux.Vars(r)
		userId := params["userId"]
		var user models.User
		defer cancel()
		println("Reached-1")
		objectId, _ := primitive.ObjectIDFromHex(userId)
		println("Reached-2")
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
		}
		println("Reached-3")
		if validationErr := validate.Struct(&user); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := responses.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
		println("Reached-4")
		res, err := userCollection.UpdateOne(ctx, bson.M{"id": objectId}, bson.M{"$set": update})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		println("Reached-5")
		var updatedUser models.User
		if res.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objectId}).Decode(&updatedUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		println("Reached-6")
		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteAUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		params := mux.Vars(r)
		userId := params["userId"]
		defer cancel()

		objectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fmt.Println("Error while getting object id", err)
			return
		}
		res, err := userCollection.DeleteOne(ctx, bson.M{"id": objectId})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		if res.DeletedCount < 1 {
			w.WriteHeader(http.StatusNotFound)
			response := responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User deleted successfully!"}}
		json.NewEncoder(w).Encode(response)
	}
}

func GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		var users []models.User
		defer cancel()

		res, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		defer res.Close(ctx)
		for res.Next(ctx) {
			var user models.User
			if err = res.Decode(&user); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}
			users = append(users, user)
		}
		w.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(w).Encode(response)
	}
}
