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
)

//AuthenticationMiddleware ...
type AuthenticationMiddleware struct {
	signupRequest bool
	token         string
}

//Request ...
type Request struct {
	http.Request
	user auth.User
}

//Middleware ...
func (middleware *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Preflight
		// utils.SendPreflightResponse(w, r)

		fmt.Println(r.URL.Path)
		fmt.Println()
		if r.URL.Path == auth.SignUpPath || r.URL.Path == auth.LoginPath || r.URL.Path == "/connect" || r.Method == "OPTIONS" {

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

			// auth.LoggedUser = user

			ctx := context.WithValue(r.Context(), auth.UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
			fmt.Println()

		}
	})
}

func main() {

	hub := &server.Hub{
		Clients:   make(map[string]*server.Client),
		DoServe:   make(chan *server.Client),
		Broadcast: make(chan server.Message),
		Entered:   make(chan *server.Client),
	}

	r := mux.NewRouter()

	go hub.Run()

	auth.InitAuthDb()
	hub.InitServer()

	// preflight conditions check
	r.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.SendPreflightResponse(w, r)
		})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})

	r.HandleFunc("/connect", auth.ConnectHandler).Methods("POST")
	r.HandleFunc("/signup", auth.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ConnectionHandler(hub, w, r)
	})

	r.HandleFunc("/findfriends", server.GetFriendsListHandler).Methods("POST")

	authMiddlerware := AuthenticationMiddleware{}

	r.Use(authMiddlerware.Middleware)

	//Server Listen...
	// port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":8080", r))

}

// Segregate Msgs for each user in ui
// make authLoggedUser stateless.
