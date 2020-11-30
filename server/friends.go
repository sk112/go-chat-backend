package server

import (
	"chat/auth"
	"chat/db"
	"chat/utils"
	"fmt"
	"net/http"
)

//FriendsList ...
type FriendsList struct {
	Firstname string
	Lastname  string
	UserID    string `json:"userid"`
}

//GetFriendsListHandler ...
func GetFriendsListHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(auth.UserKey).(auth.User).UserID

	fmt.Println("userid", userid)
	var users []FriendsList

	db.DB.Table("users").Model(&FriendsList{}).Where("user_id <> ?", userid).Find(&users)

	utils.Send(w, "friends: friends list", http.StatusOK, " List of Friends", struct {
		Users []FriendsList
	}{
		Users: users,
	})
}
