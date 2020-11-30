package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	//DBPath for db init
	DBPath = "C:\\Users\\mohammak\\Documents\\projects-learning\\go-learning\\chat\\backend\\db.sqlite3"
)

//DB ...
var DB *gorm.DB

func init() {

	var err error
	DB, err = gorm.Open(sqlite.Open(DBPath), &gorm.Config{})

	if err != nil {
		log.Fatal("DB init failed")
	}

}
