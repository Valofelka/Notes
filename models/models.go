package models

import "time"

type Note struct {
	Text      string    `json:"text" csv:"text"`
	Id        int       `json:"id" csv:"id"`
	Title     string    `json:"title" csv:"title"`
	CreatedAt time.Time `json:"createdAt" csv:"createdAt"`
}
