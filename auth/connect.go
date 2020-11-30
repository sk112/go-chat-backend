package auth

import (
	"chat/utils"
	"fmt"
	"net/http"
)

//ConnectHandler ...
func ConnectHandler(w http.ResponseWriter, r *http.Request) {

	userid, err := utils.VerifyToken(r)

	if err != nil {
		// utils.Send(w, "connect: token error", http.StatusUnauthorized, "Token Invalid", nil)
		fmt.Println("connect: Token Error", err)
		return
	}

	utils.Send(w, "connect: successful", http.StatusOK, "Connect Succesful", struct {
		UserID string
	}{
		UserID: userid.(string),
	})
}
