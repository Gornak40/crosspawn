package models

import "gorm.io/gorm"

type Contest struct {
	gorm.Model

	EjudgeID uint `gorm:"unique;not null"`

	ReviewActive bool `gorm:"not null"`
}
