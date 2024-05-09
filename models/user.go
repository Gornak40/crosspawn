package models

import (
	"gorm.io/gorm"
)

// TODO: set EjudgeID and hash password.
type User struct {
	gorm.Model

	EjudgeID       uint   `gorm:"not null"`
	EjudgeLogin    string `gorm:"unique;not null;type:varchar(128)"`
	EjudgePassword string `gorm:"not null;type:varchar(128)"`

	ReviewAproveCount uint `gorm:"not null"`
	ReviewRejectCount uint `gorm:"not null"`
}

func NewUserFromForm(login, password string) *User {
	return &User{
		EjudgeLogin:    login,
		EjudgePassword: password,
	}
}
