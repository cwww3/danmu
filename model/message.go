package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	User string `json:"user"`
	Content string `json:"content"`
	Room string `json:"room"`
}
