package entity

import "time"

type URL struct {
	ID         string     `json:"id" bson:"_id"`
	Original   string     `json:"url" bson:"url"`
	CreatedAt  time.Time  `json:"created_at" bson:"created_at"`
	Hits       int64      `json:"hits,omitempty" bson:"hits,omitempty"`
	LastAccess *time.Time `json:"last_access,omitempty" bson:"last_access,omitempty"`
}

type URLDTO struct {
	URL string `json:"url"`
}
