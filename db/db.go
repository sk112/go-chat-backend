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

	dsn := "host=ec2-54-163-47-62.compute-1.amazonaws.com user=bfrsdcmkvfmiwf password=b2672c03597978511affde364b3f8cb3c64944228627c40d4ed6318a581595c8 dbname=d4urqcaecbfcmp port=5432"

	var err error
	// DB, err = gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB init failed")
	}

}
