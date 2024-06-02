package main

import (
	"fmt"
	"html"
	"net/http"
)

const PORT string = "8080"

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello, %q", html.EscapeString(req.URL.Path))
	})

	http.HandleFunc("/hi", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hi")
	})

	fmt.Println("Server started on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
