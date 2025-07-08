package entity

import (
	"encoding/json"
	"time"
)

type Event interface {
	Bytes() []byte
	SetData([]byte)
}

type AccessEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewAccessEvent(data []byte) *AccessEvent {
	obj := AccessEvent{}
	_ = json.Unmarshal(data, &obj)

	return &obj
}

func (ae *AccessEvent) Bytes() []byte {
	data, _ := json.Marshal(ae)
	return data
}

func (ae *AccessEvent) SetData(data []byte) {
	_ = json.Unmarshal(data, ae)
}
