package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model

	Title string `gorm:"size:30"`
	Text  string `gorm:"size:300"`
}
