package server

import (
	"chat/auth"
	"chat/db"
	"fmt"
)

//InitServer ...
func (h *Hub) InitServer() {

	var users []auth.User

	db.DB.Table("users").Find(&users)

	for _, user := range users {
		client := &Client{
			hub:      h,
			conn:     nil,
			user:     user,
			message:  make(chan []Message),
			queueMsg: make([]Message, 0),
		}

		client.hub.DoServe <- client
	}

	fmt.Println("init server: clients created and serving")
}
