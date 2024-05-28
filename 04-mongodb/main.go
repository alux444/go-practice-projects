package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alux444/go-mongodb/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PORT = "8081"
const hostName = "localhost"
const mongoPort = "27017"

func getSession() *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s", hostName, mongoPort)
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func main() {
	router := httprouter.New()
	client := getSession()
	userController := controllers.NewUserController(client)
	router.GET("/user/:id", userController.GetUser)
	router.POST("/user", userController.CreateUser)
	router.DELETE("/user/:id", userController.DeleteUser)
	fmt.Println("Server started at port " + PORT)
	log.Fatal(http.ListenAndServe("localhost:"+PORT, router))
}
