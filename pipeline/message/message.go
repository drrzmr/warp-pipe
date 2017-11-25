package message

import "time"

// Payload type
type Payload map[string]interface{}

// Message struct
type Message struct {
	timestamp time.Time
	payload   Payload
}

// New return new message with given payload
func New(payload Payload) Message {
	return Message{
		timestamp: time.Now(),
		payload:   payload,
	}
}

// Get return payload for given key
func (m *Message) Get(key string) interface{} {
	return m.payload[key]
}

// Timestamp return message creation timestamp
func (m *Message) Timestamp() time.Time {
	return m.timestamp
}
