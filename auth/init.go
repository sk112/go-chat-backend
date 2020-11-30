package auth

import (
	"chat/db"
	"fmt"
	"log"
)

//UserCtx ...
type UserCtx string

//UserKey ...
var UserKey = UserCtx("user")

//InitAuthDb ...
func InitAuthDb() {
	fmt.Println(" auth init: ")

	err := db.DB.AutoMigrate(&User{})

	if err != nil {
		log.Fatal("Migration failed")
	}
}
