package main

import (
	"chat/auth"
	"chat/db"
	"chat/server"
	"chat/utils"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
			userid, err := utils.VerifyToken(r)

			user := auth.User{}
			err = db.DB.Table("users").Where("user_id = ?", userid).First(&user).Error

			if err != nil {
				http.Error(w, "Auth Failed/ User Does not exists", http.StatusBadRequest)
				return
			}

			log.Println("middleware: userid - token in middleware :", userid)

			ctx := context.WithValue(r.Context(), auth.UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
			fmt.Println()

		}
	})
}

func main() {

	hub := server.Hub{
		Pool: make(map[string]*server.Client),
	}
	hubconns := server.HubConnections{
		Conns:       make(map[string]*websocket.Conn),
		PublishUser: make(chan auth.User),
	}

	r := mux.NewRouter()

	auth.InitAuthDb()
	go hub.RunPool()

	// preflight conditions check
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.SendPreflightResponse(w, r)
		})

	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/ws", hub.WSHandler)
	r.HandleFunc("/joinhub", hubconns.JoinHubHandler)
	r.HandleFunc("/users", server.GetFriendsHandler).Methods("GET")
	authMiddlerware := AuthenticationMiddleware{}

	r.Use(authMiddlerware.Middleware)

	//Server Listen...
	// port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":8080", r))

}
