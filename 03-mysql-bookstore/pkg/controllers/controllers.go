package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alux444/go-bookstore/pkg/models"
	"github.com/alux444/go-bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

var newBook models.Book

func GetBooks(res http.ResponseWriter, req *http.Request) {
	allBooks := models.GetAllBooks()
	result, _ := json.Marshal(allBooks)
	res.Header().Set("Content-Type", "pkglication/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)
}

func GetBookById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["bookId"]
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing ")
	}
	thisBook, _ := models.GetBookById(bookId)
	result, _ := json.Marshal(thisBook)
	res.Header().Set("Content-Type", "pkglication/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)
}

func DeleteBook(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["bookId"]
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing ")
	}
	thisBook := models.DeleteBookById(bookId)
	result, _ := json.Marshal(thisBook)
	res.Header().Set("Content-Type", "pkglication/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)
}

func CreateBook(res http.ResponseWriter, req *http.Request) {
	CreateBook := &models.Book{}
	utils.ParseBody(req, CreateBook)
	b := CreateBook.CreateBook()
	result, _ := json.Marshal(b)
	res.WriteHeader(http.StatusOK)
	res.Write(result)
}

func UpdateBook(res http.ResponseWriter, req *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(req, updateBook)
	vars := mux.Vars(req)
	id := vars["bookId"]
	bookId, err := strconv.ParseInt(id, 0, 0)
	if err != nil {
		fmt.Println("Error parsing body")
	}
	bookDetails, db := models.GetBookById(bookId)

	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	db.Save(&bookDetails)
	result, _ := json.Marshal(bookDetails)
	res.Header().Set("Content-Type", "pkglication/json")
	res.WriteHeader(http.StatusOK)
	res.Write(result)
}
