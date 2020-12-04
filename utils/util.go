package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//GetJSONPOSTBody parse through json post body and return body string.
func GetJSONPOSTBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		return []byte(""), err
	}

	if err := r.Body.Close(); err != nil {
		return []byte(""), err
	}

	return body, nil
}

//ExtractToken ...
func ExtractToken(r *http.Request) string {

	var bearToken string
	if r.URL.Path == "/ws" || r.URL.Path == "/joinhub" {
		bearToken = r.FormValue("token")
	} else {
		bearToken = r.Header.Get("Authorization")

		//normally Authorization the_token_xxx
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			return strArr[1]
		}
	}

	return bearToken
}

//VerifyToken ...
func VerifyToken(r *http.Request) (interface{}, error) {
	tokenString := ExtractToken(r)

	// fmt.Println("tokenString", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("SECRET_KEY"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok == false {
		log.Fatal(err)
	}

	return claims["userid"], nil
}
