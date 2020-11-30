package server

import (
	"chat/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//Client ...
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	user     auth.User
	message  chan []Message
	queueMsg []Message
}

//Message ...
type Message struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TextMessage string `json:"message"`
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *Client) writePump(r *http.Request) {

	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.message:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("next writer error-unexpected close : %v", err)
					return
				}
				fmt.Println("next writer: error", err)
				return
			}

			// fmt.Println("writePump: ", c.user.UserID, " ----should be-->", msg.From, "---logged in user--->", r.Context().Value(auth.UserKey).(auth.User))

			marshalledMsg, err := json.Marshal(&msg)

			fmt.Println("writePump :", string(marshalledMsg))
			w.Write(marshalledMsg)

			w.Close()
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		// c.conn.Close()
		c.conn = nil
	}()

	for {

		msg := Message{}

		_, t, err := c.conn.ReadMessage()

		// if err != nil {
		// 	fmt.Println("read pump: read failed: err : ", err)
		// 	return
		// }

		fmt.Println("read pump: message", msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("read pump: error - unexpected close: %v", err)
				return
			}
			fmt.Println("read pump: read failed: err : ", err)
			return
		}

		err = json.Unmarshal(t, &msg)

		fmt.Println(msg, err)
		c.hub.Broadcast <- msg
	}
}

//ConnectionHandler ...
func ConnectionHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("upgrader error", err)
		return
	}

	client := hub.Clients[r.Context().Value(auth.UserKey).(auth.User).UserID]

	client.conn = conn

	go client.readPump()
	go client.writePump(r)

	client.hub.Entered <- client
}
