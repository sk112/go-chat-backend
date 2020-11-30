package auth

import (
	"chat/db"
	"chat/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//SignUpForm ...
type SignUpForm struct {
	Firstname string
	Lastname  string
	ID        string `json:"userid"`
	Password1 string
	Password2 string
}

// SignUpHandler ...
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In wsign up handler")

	form := SignUpForm{}

	body, err := utils.GetJSONPOSTBody(r)

	fmt.Println(string(body))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(body, &form)

	if err != nil {
		log.Fatal("Json Unmarshaling failed!")
		return
	}

	var count int64

	db.DB.Table("users").Where("user_id = ?", form.ID).Count(&count)

	if count > 0 {
		utils.Send(w, "sign up: User name already taken", http.StatusBadRequest, "User name already taken", nil)
		return
	}

	if form.ID == "" || form.Password1 == "" || form.Password2 == "" {
		utils.Send(w, "sign up: required field is empty", http.StatusBadRequest, "(One of) Required Field is null", nil)
		return
	}

	if form.Password1 != form.Password2 {
		utils.Send(w, "sign up: Passwords does not match", http.StatusBadRequest, "Passwords does not match", nil)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password1), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("id:", form.ID)
	fmt.Println("password:", form.Password1)
	fmt.Println("password:", form.Password2)

	user := User{
		UserID:    form.ID,
		Password:  string(hash),
		Firstname: form.Firstname,
		Lastname:  form.Lastname,
	}

	err = db.DB.Create(&user).Error

	if err != nil {

		utils.Send(w, "sign up: DB Error", http.StatusBadRequest, "Db Error", nil)
		return

	}

	// w.Write([]byte("Registered Succesfully"))
	utils.Send(w, "sign up: Registered Succesfully", http.StatusOK, "Registered Succesfully", nil)
}
