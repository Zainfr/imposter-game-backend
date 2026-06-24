package hub

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Zainfr/imposter-game-backend/internal/game"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan IncomingMessage
	register chan *Client
	unregister chan *Client
	game	*game.Game
}

func NewHub() *Hub{
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan IncomingMessage),
		register: make(chan *Client),
		unregister: make(chan *Client),
		game: &game.Game{
			Players: make(map[string]*game.Player),
			CurrentPhase: game.PhaseLobby,
		},
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
				switch msg.Type {
				case "join_room":
					var p struct { PlayerName string `json:"player_name"`}
					if err := json.Unmarshal(msg.Payload, &p); err != nil{
						log.Printf("invalid join_room payload: %v", err)
						continue
					}
					id := generateID()
					h.game.Players[id] = &game.Player{ ID: id, Name: p.PlayerName}
					msg.Client.ID = id
					
					broadcastMsg := OutgoingMessage{
						Type: "player_joined",
						Payload: PlayerJoinedPayload{ID: id, Name: p.PlayerName},
					}
					h.broadcastToAll(broadcastMsg)
				case "start_game": 
					if err := h.game.StartGame(game.WordList); err != nil {
						log.Println("start game error:", err)
						continue
					}
					for client := range h.clients {
						word := h.game.WordForPlayer(client.ID)
						msg := OutgoingMessage{
							Type: "game_started",
							Payload: map[string]interface{}{
								"secret_word": word,
							},
						}
						client.send <- msg
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

func (h *Hub) broadcastToAll(msg OutgoingMessage){
	for client := range h.clients {
		select{
			case client.send <- msg:
			default:
				close(client.send)
				delete(h.clients, client)
		}
	}
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &Client{
		hub: h,
		conn: conn,
		send: make(chan OutgoingMessage, 256),
	}
	h.register <- client
	go client.readPump()
	go client.writePump()

	return nil
}