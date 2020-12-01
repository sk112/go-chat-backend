package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	//DBPath for db init
	DBPath = "C:\\Users\\mohammak\\Documents\\projects-learning\\go-learning\\chat\\backend\\db.sqlite3"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "calhounio_demo"
)

//DB ...
var DB *gorm.DB

func init() {

	dsn := "host=ec2-52-71-153-228.compute-1.amazonaws.com user=xelxqkwnnrmwmi password=202ace837067b1fcb5d81a167019d1e1e6e1e033a1aa6cbab33562ad225bcef1 dbname=da6lspqu3sdmga port=5432"

	var err error
	// DB, err = gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB init failed")
	}

}
