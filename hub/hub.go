package hub

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func NewHub() *Hub{
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select{
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok{
					delete(h.clients, client)
					close(client.send)
				}
			case msg := <-h.broadcast:
				for client := range h.clients{
					select {
						case client.send <- msg:
						default:
							close(client.send)
							delete(h.clients, client)
					}
				}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{
		hub: h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.register <- client
	go client.readPump()
	go client.writePump()

	return nil
}