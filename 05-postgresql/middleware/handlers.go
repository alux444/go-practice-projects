package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type response struct {
	Id      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const PSQL_HOST string = "localhost"
const PSQL_PORT string = "5432"
const PSQL_DB string = "testing"
const PSQL_USER string = "yourusername"
const PSQL_PASS string = "yourpassword"

const URI string = "postgres://" + PSQL_USER + ":" + PSQL_PASS + "@" + PSQL_HOST + ":" + PSQL_PORT + "/" + PSQL_DB + "?sslmode=disable"

func createConnection() *sql.DB {
	db, err := sql.Open("postgres", URI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to PSQL db")
	return db
}

func GetStock(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id param: %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Error getting stock: %v", err)
	}
	json.NewEncoder(res).Encode(stock)
}

func GetStocks(res http.ResponseWriter, req *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("Error getting all stocks: %v", err)
	}
	json.NewEncoder(res).Encode(stocks)
}

func CreateStock(res http.ResponseWriter, req *http.Request) {
	var stock models.Stock
	err := json.NewDecoder(req.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the body: %v", err)
	}
	insertId := insertStock(stock)

	response := response{
		Id:      insertId,
		Message: "Stock created successfuly",
	}

	json.NewEncoder(res).Encode(response)
}

func UpdateStock(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Error converting ID parameter: %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(req.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Error decoding request body: %v", err)
	}
	rows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated, affected %v rows", rows)
	response := response{
		Id:      int64(id),
		Message: msg,
	}
	json.NewEncoder(res).Encode(response)
}

func DeleteStock(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Error converting id: %v", err)
	}
	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stock deleted, affected %v rows", deletedRows)
	response := response{
		Id:      int64(id),
		Message: msg,
	}
	json.NewEncoder(res).Encode(response)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	query := `INSERT INTO stocks(name,price,company) VALUES ($1,$2,$3) RETURNING stockid`
	var id int64
	err := db.QueryRow(query, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute query %v", err)
	}

	fmt.Printf("Inserted stock with id: %v\n", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stock models.Stock
	query := `SELECT * FROM stocks WHERE stockid=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Printf("No rows returned for id: %v\n", id)
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan rows for id: %v", err)
	}

	return stock, err
}

func getAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()
	var stocks []models.Stock
	query := `SELECT * FROM stocks`
	rows, err := db.Query(query)

	if err != nil {
		log.Fatalf("Error in query")
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Error scanning rows in stock")
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	query := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err := db.Exec(query, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Error setting stocks: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error in checking rows affected: %v", err)
	}
	fmt.Printf("Total rows affected: %v\n", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()
	query := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(query, id)
	if err != nil {
		log.Fatalf("Error deleting stocks: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error in checking rows affected: %v", err)
	}
	fmt.Printf("Total rows affected: %v\n", rowsAffected)
	return rowsAffected
}
