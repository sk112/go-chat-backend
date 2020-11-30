package server

import (
	"log"
)

//Hub ...
type Hub struct {
	Clients   map[string]*Client
	DoServe   chan *Client
	Entered   chan *Client
	Broadcast chan Message
}

//Run ...
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.DoServe:
			h.Clients[client.user.UserID] = client
		case client := <-h.Entered:
			log.Println("hub: user enter >", client.user.UserID)

			if len(client.queueMsg) > 0 {
				client.message <- client.queueMsg
			}
			// h.Clients[client.user.UserID] = client
		case msg := <-h.Broadcast:
			m := []Message{msg}

			if h.Clients[msg.To].conn != nil {

				h.Clients[msg.To].message <- m
			} else {
				t := append(h.Clients[msg.To].queueMsg, msg)
				h.Clients[msg.To].queueMsg = t
			}

			if h.Clients[msg.From] != nil {
				h.Clients[msg.From].message <- m
			}
			// else {
			// 	t := append(h.Clients[msg.From].queueMsg, msg)
			// 	h.Clients[msg.From].queueMsg = t
			// }

		}
	}
}
