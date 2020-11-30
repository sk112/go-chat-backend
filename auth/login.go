package auth

import (
	"chat/db"
	"chat/utils"
	"log"
	"time"

	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//LoginForm ...
type LoginForm struct {
	UserID   string `json:"userid"`
	Password string
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
		// fmt.Println("User Does not exists")
		// http.Error(w, err.Error(), http.StatusBadRequest)
		utils.Send(w, "auth: user extraction failed", http.StatusBadRequest, "user does not exists", []byte("{}"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		// http.Error(w, " Password does not match", http.StatusBadRequest)
		utils.Send(w, "auth: password does not match", http.StatusBadRequest, "password Incorrect", []byte("{}"))
		return
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
