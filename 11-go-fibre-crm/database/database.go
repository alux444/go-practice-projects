package database

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var (
	Connection *gorm.DB
)
