package server

import (
	"chat/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//ClientInHub :
type ClientInHub struct {
	Conn       *websocket.Conn
	User       string
	TellClient chan MessageToClient
	Room       *HubRoom
	QueueIdx   int
}

//MessageToClient :
type MessageToClient struct {

	//Entering = 1
	//Entered = 2
	//Leaving = 3
	Action  int
	ClientS interface{}
}

//HubRoom ...
type HubRoom struct {
	Conns   map[string]*ClientInHub
	Entered chan *ClientInHub
	Left    chan *ClientInHub
}

//JoinHubHandler :
func (h *HubRoom) JoinHubHandler(w http.ResponseWriter, r *http.Request) {

	try := r.FormValue("try")

	fmt.Println(">>>>>>>>>>Retrying count :", try)
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

	userid := r.Context().Value(auth.UserKey).(auth.User).UserID

	clientInHub := &ClientInHub{
		Conn:       conn,
		User:       userid,
		TellClient: make(chan MessageToClient),
		Room:       h,
	}

	go clientInHub.readFromUser(conn, userid, r)
	go clientInHub.WriteToUser()

	h.Entered <- clientInHub
}

func (cih *ClientInHub) readFromUser(conn *websocket.Conn, userid string, r *http.Request) {
	// go cih.waitTillLeft()
	defer func() {
		cih.Conn.Close()
	}()

	cih.Conn.SetReadDeadline(time.Now().Add(24 * time.Hour))

	defer func() {
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("<error>  ReadMessage error-unexpected close : %v", err)

				cih.Room.Left <- cih
				return
			}
			fmt.Printf("<error>  ReadMessage error-unexpected close : %v", err)

			cih.Room.Left <- cih
			return
		}
	}
}

//WriteToUser :
func (cih *ClientInHub) WriteToUser() {
	defer func() {
		cih.Conn.Close()
	}()

	for {
		select {
		case msg := <-cih.TellClient:

			// for k, v := range h.Conns {
			cih.Conn.SetWriteDeadline(time.Now().Add(24 * time.Hour))
			w, err := cih.Conn.NextWriter(websocket.TextMessage)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("<error>  NextWriter error-unexpected close : %v", err)

					cih.Room.Left <- cih
					return
				}
				fmt.Printf("<error>  NextWriter error-unexpected close : %v", err)

				cih.Room.Left <- cih
				return
			}

			var jsonUserIDs []byte

			if msg.Action == 1 || msg.Action == 3 {
				jsonUserIDs, err = json.Marshal(&struct {
					ActionType int
					UserID     string
				}{
					ActionType: msg.Action,
					UserID:     msg.ClientS.(string),
				})

				if err != nil {
					fmt.Println("<error> write error >><>", err)
				}
			} else if msg.Action == 2 {
				jsonUserIDs, err = json.Marshal(&struct {
					ActionType int
					UserID     []string
				}{
					ActionType: msg.Action,
					UserID:     msg.ClientS.([]string),
				})

				if err != nil {
					fmt.Println("<error> write error >><>", err)
				}
			}

			w.Write(jsonUserIDs)
			w.Close()
			// }

		}
	}
}

//WaitTillLeft :
func (h *HubRoom) WaitTillLeft() {
	for {
		select {
		case clientInHub := <-h.Left:
			fmt.Println("<log> client left : ", clientInHub.User)
			clientInHub.Conn.Close()
			delete(h.Conns, clientInHub.User)

			var users []string

			for _, v := range h.Conns {
				users = append(users, v.User)
			}

			for _, user := range users {
				h.Conns[user].TellClient <- MessageToClient{Action: 3, ClientS: clientInHub.User}
			}
		}
	}
}

//EnterRoom :
func (h *HubRoom) EnterRoom() {

	for {
		select {
		case clientInHub := <-h.Entered:

			var users []string

			for _, v := range h.Conns {
				users = append(users, v.User)
			}
			for _, user := range users {
				h.Conns[user].TellClient <- MessageToClient{Action: 1, ClientS: clientInHub.User}
			}

			h.Conns[clientInHub.User] = clientInHub

			if len(users) != 0 {
				clientInHub.TellClient <- MessageToClient{Action: 2, ClientS: users}
			}

		}

	}
}
