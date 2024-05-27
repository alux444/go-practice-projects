package main

import (
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
	log.Fatal(http.ListenAndServe("localhost:"+PORT, router))
}
