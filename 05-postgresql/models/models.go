package models

type Stock struct {
	StockId int64  `json:"stockId"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Company string `json:"company"`
}
