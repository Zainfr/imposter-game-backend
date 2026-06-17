package hub

import (
	"encoding/json"
	"fmt"
)

type IncomingMessage struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
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