package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alux444/go-mongodb/models"
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

func (uc UserController) GetUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.client.Database("mongo-golang").Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}

func (uc UserController) CreateUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	user := models.User{}
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	user.Id = primitive.NewObjectID()
	collection := uc.client.Database("mongo-golang").Collection("users")
	collection.InsertOne(context.Background(), user)
	userJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s\n", userJson)
}

func (uc UserController) DeleteUser(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := uc.client.Database("mongo-golang").Collection("users")
	result, err := collection.DeleteOne(context, bson.M{"_id": oid})
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if result.DeletedCount == 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "Deleted user %s\n", oid)
}
