package models

type Stock struct {
	StockId int64   `json:"stockId"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Company string  `json:"company"`
}
