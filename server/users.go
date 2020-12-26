package server

import (
	"chat/auth"
	"chat/utils"
	"net/http"
)

//GetFriendsHandler ...
func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {

	var users []auth.User

	// db.DB.Table("users").Find(&users)

	// jsonUsers, err := json.Marshal(&users)

	// if err != nil {
	// 	fmt.Println("<error> get friends list : json marshal error:,", err)
	// 	return
	// }

	utils.Send(w, "get users list", http.StatusOK, "users list", users)
}
