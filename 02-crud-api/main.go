package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	id       string    `json:"id"`
	isbn     string    `json:"isbn"`
	title    string    `json:"title"`
	director *Director `json:"director"`
}

type Director struct {
	firstName string `json:"firstName"`
	lastName  string `json:"lastName"`
}

var movies []Movie

const PORT string = "8080"

func getMovies(res http.ResponseWriter, req *http.Request) {

}

func getMovie(res http.ResponseWriter, req *http.Request) {

}

func createMovie(res http.ResponseWriter, req *http.Request) {

}

func updateMovie(res http.ResponseWriter, req *http.Request) {

}

func deleteMovie(res http.ResponseWriter, req *http.Request) {

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Server started at port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
