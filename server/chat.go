package server

import (
	"chat/auth"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//Hub :
type Hub struct {
	// Users[test][test1] = ChatUSer
	Users map[string]map[string]*ChatUser
}

//ChatUser :
type ChatUser struct {
	From    string
	With    string
	Conn    *websocket.Conn
	Hub     *Hub
	Message chan []byte
}

func (c *ChatUser) readMessage() {
	defer func() {
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(24 * time.Hour))
	for {

		_, msg, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("<error>  ReadMessage error-unexpected close : %v", err)

				return
			}
			fmt.Printf("<error>  ReadMessage error-unexpected close : %v", err)

			return
		}

		fmt.Println("message: ", c.With, c.From, string(msg))
		//Forwarding message by c.with->c.from connection.
		c.Hub.Users[c.With][c.From].Message <- msg

	}
}

func (c *ChatUser) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg := <-c.Message:
			c.Conn.SetWriteDeadline(time.Now().Add(24 * time.Hour))
			w, err := c.Conn.NextWriter(websocket.TextMessage)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("<error>  NextWriter error-unexpected close : %v", err)
					return
				}
				fmt.Printf("<error>  NextWriter error-unexpected close : %v", err)

				return
			}

			w.Write(msg)
			w.Close()
		}
	}
}

//ConnectHandler :
func (h *Hub) ConnectHandler(w http.ResponseWriter, r *http.Request) {

	with := r.FormValue("with")

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

	chat := &ChatUser{
		From:    userid,
		With:    with,
		Conn:    conn,
		Hub:     h,
		Message: make(chan []byte),
	}

	if h.Users[userid] == nil {
		h.Users[userid] = make(map[string]*ChatUser)
	}

	h.Users[userid][with] = chat

	go chat.readMessage()
	go chat.writeMessage()
}
