package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var movies []Movie

const PORT string = "8080"

func getMovies(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(movies)
}

func getMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(res).Encode(item)
			return
		}
	}
}

func deleteMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range movies {
		if item.Id == params["id"] {
			//delete one movie by adjusting index of slice
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(res).Encode(movies)
}

func createMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(req.Body).Decode(&movie)

	movie.Id = strconv.Itoa(rand.IntN(1000000))
	movies = append(movies, movie)
	json.NewEncoder(res).Encode(movie)
}

func updateMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(req.Body).Decode(&movie)
			movie.Id = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(res).Encode(movie)
			return
		}
	}
}

func main() {
	movie1 := Movie{Id: "1", Isbn: "1234", Title: "movie1", Director: &Director{FirstName: "young", LastName: "timmy"}}
	movie2 := Movie{Id: "2", Isbn: "1235", Title: "movie2", Director: &Director{FirstName: "old", LastName: "timmy"}}
	movies = append(movies, movie1)
	movies = append(movies, movie2)

	router := mux.NewRouter()
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Server started at port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
