package model

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	User string `json:"user"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type User struct {
	gorm.Model
	Name string `json:"name"`
}
