package server

import (
	"chat/auth"
	"chat/db"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

//HubConnections map[string]*websocket.Conn
type HubConnections struct {
	Conns       map[string]*websocket.Conn
	PublishUser chan auth.User
}

//Hub ...
type Hub struct {
	Pool    map[string]*Client
	Entered chan *Client
}

//RunPool : Init for all Client connection to listen on the channel.
func (h *Hub) RunPool() {

	for {
		select {
		case registerClient := <-h.Entered:
			h.Pool[registerClient.ConnID] = registerClient

			for k, v := range h.Pool {
				fmt.Println(k, v)
			}
		}
	}
}

func (h *Hub) writeToHub(r *http.Request) {
	for {

		select {
		case user, ok := <-h.Entered:

			fmt.Println(user, ok)
		}

	}
}

//JoinHubHandler ...
func (h *HubConnections) JoinHubHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("<error> Join hub : upgrader error", err)
	}

	h.Conns[r.Context().Value(auth.UserKey).(auth.User).UserID] = conn

	go h.BroadCastToUser()

	h.PublishUser <- r.Context().Value(auth.UserKey).(auth.User)
}

//BroadCastToUser ...
func (h *HubConnections) BroadCastToUser() {

	for {
		select {
		case user, ok := <-h.PublishUser:

			if !ok {
				fmt.Println("<error> write to user: received an error")
			}

			for k, v := range h.Conns {

				if v != nil {
					w, err := v.NextWriter(websocket.TextMessage)

					if err != nil {

						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							fmt.Printf("<error> next writer error-unexpected close : %v", err)

							db.DB.Table("users").Where("user_id =", k).Delete(&auth.User{})
							return
						}
						fmt.Println("<error> write to user: next wrtiter failed: ", err)
						return
					}

					w.Write([]byte(user.UserID))
				}

			}
		}
	}
}

func (h *HubConnections) ReadFromUser() {

	for k, v := range h.Conns{
		go func(){
			r ,err := v.ReadMessage()

			if err != nil{
				fmt.Println("")
			}
		}

	}
}
