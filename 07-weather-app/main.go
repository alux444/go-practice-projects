package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	ApiKey string `json:"apiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

const SERVER_PORT = "8080"

func hello(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("hello\n"))
}

func query(city string) (weatherData, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading ENV: %v\n", err)
	}
	key := os.Getenv("API_KEY")
	res, err := http.Get("https://api.openweathermap.org/data/2.5/weather?APPID=" + key + "&q=" + city)

	if err != nil {
		log.Printf("%v\n", err)
		return weatherData{}, err
	}
	defer res.Body.Close()

	var data weatherData
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		log.Printf("%v\n", err)
		return weatherData{}, err
	}

	return data, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(res http.ResponseWriter, req *http.Request) {
			city := strings.SplitN(req.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(data)
		})
	fmt.Println("Server started on port: " + SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, nil)
}
