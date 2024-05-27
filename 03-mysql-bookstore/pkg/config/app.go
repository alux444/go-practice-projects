package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open("mysql", "test:Password!12345@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	db = d
}

func GetDb() *gorm.DB {
	return db
}
