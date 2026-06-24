package hub

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct{
	hub *Hub
	conn *websocket.Conn
	send chan OutgoingMessage
	ID string
}

func (c *Client) readPump(){
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, raw, err := c.conn.ReadMessage()
		if err!= nil {
			log.Println("read error:", err)
			break
		}
		msg, err2 := DecodeEnvelope(raw)
		if err2!= nil {
			log.Println("read error:", err2)
			continue
		}
		msg.Client = c
		c.hub.broadcast <- msg
	}

}

func (c *Client) writePump(){
	defer c.conn.Close()

	for msg := range c.send {
		raw, err := json.Marshal(msg)
		if err != nil {
			log.Println("marshal error:", err)
			continue 
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, raw); err != nil {
			log.Println("write error", err)
			return	
		}
	}
}