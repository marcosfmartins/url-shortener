package entity

import "time"

type StatisticDTO struct {
	Hits       int64      `json:"hits" bson:"hits,omitempty"`
	LastAccess *time.Time `json:"lastAccess,omitempty" bson:"lastAccess,omitempty"`
}

func URLtoStatistic(url URL) StatisticDTO {
	return StatisticDTO{
		Hits:       url.Hits,
		LastAccess: url.LastAccess,
	}
}
