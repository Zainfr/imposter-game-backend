package room

import (
	"sync"

	"github.com/Zainfr/imposter-game-backend/internal/hub"
)

type Registry struct {
	mu sync.Mutex
	rooms map[string]*hub.Hub
}

func NewRegistry() *Registry {
	return &Registry{
		rooms: make(map[string]*hub.Hub),
	}
}

func (r *Registry) CreateRoom(code string) *hub.Hub{
	r.mu.Lock()
	defer r.mu.Unlock()

	h := hub.NewHub()
	go h.Run()

	r.rooms[code] = h
	return h
}

func (r *Registry) GetRoom(code string) (*hub.Hub, bool){
	r.mu.Lock()
	defer r.mu.Unlock()

	h, ok := r.rooms[code]
	return h, ok
}
