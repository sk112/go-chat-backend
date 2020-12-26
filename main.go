package main

import (
	"chat/auth"
	"chat/server"
	"chat/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//AuthenticationMiddleware ...
type AuthenticationMiddleware struct {
}

//Middleware ...
func (middleware *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		if r.URL.Path == auth.LoginPath || r.Method == "OPTIONS" {

			next.ServeHTTP(w, r)
		} else {
			userid, _ := utils.VerifyToken(r)

			user := auth.User{
				UserID: userid.(string),
			}

			log.Println("middleware: userid - token in middleware :", userid)

			ctx := context.WithValue(r.Context(), auth.UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
	})
}

func main() {

	hubRoom := &server.HubRoom{
		Conns:   make(map[string]*server.ClientInHub),
		Entered: make(chan *server.ClientInHub),
		Left:    make(chan *server.ClientInHub),
	}

	r := mux.NewRouter()
	// go hubRoom.WriteToUsers()
	go hubRoom.EnterRoom()
	go hubRoom.WaitTillLeft()

	// preflight conditions check
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.SendPreflightResponse(w, r)
		})

	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/joinhub", hubRoom.JoinHubHandler)

	hub := &server.Hub{
		Users: make(map[string]map[string]*server.ChatUser),
	}

	r.HandleFunc("/connect", hub.ConnectHandler)

	authMiddlerware := AuthenticationMiddleware{}

	r.Use(authMiddlerware.Middleware)

	//Server Listen...
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, r))

}
