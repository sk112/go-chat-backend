package auth

import (
	"chat/db"
	"chat/utils"
	"log"
	"time"

	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//LoginForm ...
type LoginForm struct {
	UserID string `json:"userid"`
}

//LoginHandler ...
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}

	body, err := utils.GetJSONPOSTBody(r)

	if err != nil {

		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Send(w, "json post formation failed", http.StatusBadRequest, "json post formation failed", []byte("{}"))
		return
	}

	err = json.Unmarshal(body, &form)

	if err != nil {
		log.Fatal("JSON Unmarshal Failing")
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Send(w, "json unmarshal formation failed", http.StatusBadRequest, "json unmarchaling formation failed", []byte("{}"))
		return
	}

	user := User{}
	err = db.DB.Table("users").Where("user_id = ?", form.UserID).First(&user).Error

	// fmt.Println("user err", err)
	if err != nil {
		result := db.DB.Table("users").Create(&form)

		if result.Error != nil {
			utils.Send(w, "auth: user insert failed ", http.StatusBadRequest, "user insert failed", []byte("{}"))
			return
		}

		user.UserID = form.UserID
	}

	claims := jwt.MapClaims{}

	claims["userid"] = user.UserID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte("SECRET_KEY"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Send(w, "auth: unauth", http.StatusBadRequest, "User Unauthorised", nil)
		return
	}

	log.Printf("login: token> %s\n", token)
	log.Printf("login: user> %+v\n", user)

	utils.Send(w, "auth: auth - "+user.UserID, http.StatusOK, "Login Successful", struct {
		Token  string
		UserID string
	}{
		Token:  token,
		UserID: user.UserID,
	})
}
