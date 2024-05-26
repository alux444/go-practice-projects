package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/hello" {
		http.Error(res, "404 not found.", http.StatusNotFound)
		return
	}
	if req.Method != "GET" {
		errMsg := "Method: " + req.Method + " is not supported."
		http.Error(res, errMsg, http.StatusNotFound)
		return
	}
	fmt.Fprintf(res, "Hello")
}

func formHandler(res http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Fprintf(res, "ParseForm error: %v\n ", err)
	}
	fmt.Fprintf(res, "POST request successful\n")
	fn := req.FormValue("firstName")
	ln := req.FormValue("lastName")
	fmt.Fprintf(res, "First name: %v\n", fn)
	fmt.Fprintf(res, "Last name: %v\n", ln)
}

func main() {
	port := "8080"
	//localhost:8080
	//localhost:8080/form.html
	//localhost:8080/hello
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Starting server at port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
