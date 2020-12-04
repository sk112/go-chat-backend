package server

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//Client ...
type Client struct {
	ConnID     string
	SelfConn   *websocket.Conn
	FriendConn *Client
	Message    chan []byte
}

func (c *Client) writePump() {
	for {
		select {
		case msg, ok := <-c.Message:
			fmt.Println("readpump: message writing", msg, ok)

			w, err := c.SelfConn.NextWriter(websocket.TextMessage)

			if err != nil {
				fmt.Println("write pump: next writer error.")
				return
			}

			w.Write(msg)
		}
	}
}

func (c *Client) readPump() {
	for {

		_, msg, err := c.SelfConn.ReadMessage()

		if err != nil {
			fmt.Println("read pump: read message failed", err)
			return
		}

		c.FriendConn.Message <- msg
	}
}

//AddClient ...
func (h *Hub) AddClient(friendid string, w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("<error> Add client: upgrader error", err)
	}

	client := &Client{}

	id := uuid.New().String()

	client.ConnID = id
	client.SelfConn = conn
	client.FriendConn = h.Pool[friendid]
	client.Message = make(chan []byte)

	go client.readPump()
	go client.writePump()
	// h.Entered <- client

}

//WSHandler ...
func (h *Hub) WSHandler(w http.ResponseWriter, r *http.Request) {
	fID := r.FormValue("id")

	h.AddClient(fID, w, r)
}
