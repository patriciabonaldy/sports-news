package pubsub

import (
	"time"

	"github.com/google/uuid"
)

// Message request model.
type Message struct {
	EventID   string    `json:"event_id"`
	RawData   []byte    `json:"data"`
	Timestamp time.Time `json:"at"`
}

// NewSystemMessage when we want to publish global message to all users in the room
func NewSystemMessage() (Message, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Message{}, err
	}

	return Message{
		EventID:   id.String(),
		RawData:   nil,
		Timestamp: time.Now(),
	}, nil
}
