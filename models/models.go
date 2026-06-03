package models

import (
	"time"
)

type Note struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title string `gorm:"size:30"`
	Text  string `gorm:"size:300"`
}
