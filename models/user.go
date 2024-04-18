package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
