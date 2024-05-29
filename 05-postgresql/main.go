package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
)

const SERVER_PORT string = "8080"

func main() {
	router := router.Router()
	fmt.Println("Server started on port " + SERVER_PORT)
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, router))
}
