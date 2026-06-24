package hub

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type IncomingMessage struct {
    Client *Client
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type OutgoingMessage struct {
    Type string `json:"type"`
    Payload any `json:"payload"`
}

type PlayerJoinedPayload struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type SubmitCluePayload struct {
	Clue string `json:"clue"`
}

type VotePayload struct {
	Target string `json:"target"`
}

func DecodeEnvelope(raw []byte) (IncomingMessage, error) {
    var msg IncomingMessage
    if err := json.Unmarshal(raw, &msg); err != nil {
        return IncomingMessage{}, fmt.Errorf("parsing envelope: %w", err)
    }
    if msg.Type == "" {
        return IncomingMessage{}, fmt.Errorf("missing message type")
    }
    return msg, nil
}

func generateID() string {
	return uuid.New().String()
}
