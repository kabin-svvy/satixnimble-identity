package model

import "github.com/jinzhu/gorm"

type Users struct {
	gorm.Model
	Username  string `gorm:"unique_index;not null" json:"username"`
	Email     string `gorm:"unique_index;not null" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	Firstname string `json:"firstname"`
}
