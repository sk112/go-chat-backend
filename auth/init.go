package auth

import (
	"chat/db"
	"fmt"
	"log"
)

//User ...
type User struct {
	UserID string `json:"userid"`
}

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
