package models

import "time"

type Payload struct {
	Data  string `json:"data"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type DBPayload struct {
	ID      uint      `gorm:"primaryKey"`
	Payload Payload   `gorm:"embedded;embeddedPrefix:payload_"`
	Updated time.Time `gorm:"autoUpdateTime:milli"` // Use unix milli seconds as updating time
	Created time.Time `gorm:"autoCreateTime"`
}
