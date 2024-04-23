package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	EjudgeID       uint   `gorm:"unique;not null"`
	EjudgeLogin    string `gorm:"unique;not null;type:varchar(128)"`
	EjudgePassword string `gorm:"not null;type:varchar(128)"`

	ReviewApproveCount uint `gorm:"not null"`
	ReviewRejectCount  uint `gorm:"not null"`
}
