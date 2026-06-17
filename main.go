package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/Zainfr/imposter-game-backend/internal/app"
	"github.com/Zainfr/imposter-game-backend/internal/room"
)

func generateRoomCode()	string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 6)
	rand.Read(b)
	code := make([]byte, 6)
	for i, v := range b{
		code[i] = chars[v%byte(len(chars))]
	}

	return string(code)
}

func createRoomHandler(registry *room.Registry) app.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		code := generateRoomCode()
		registry.CreateRoom(code)
		app.WriteJSON(w, http.StatusCreated, map[string]string {"code": code})
		return nil
	}
}

func wsHandler(registry *room.Registry) app.AppHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		code := r.PathValue("code")
		
		h, ok := registry.GetRoom(code)
		if !ok {
			return &app.AppError{
				Code: http.StatusNotFound,
				Message: "room not found",
				Err: fmt.Errorf("no room with code %s", code),
			}
		}

		return h.ServeWS(w,r)
	}
}

func main() {
	registry := room.NewRegistry()

	mux := http.NewServeMux()
	mux.Handle("POST /rooms", app.AppHandler(createRoomHandler(registry)))
	mux.Handle("GET /ws/{code}", app.AppHandler(wsHandler(registry)))

	http.ListenAndServe(":8080", mux)

}