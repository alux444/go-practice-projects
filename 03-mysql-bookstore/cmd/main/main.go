package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alux444/go-bookstore/pkg/routes"
	"github.com/gorilla/mux"
)

const PORT string = "8080"

func main() {
	router := mux.NewRouter()
	routes.RegisterBookstoreRoutes(router)
	http.Handle("/", router)
	fmt.Println("Server started at port: " + PORT)
	log.Fatal(http.ListenAndServe("localhost:"+PORT, router))
}
