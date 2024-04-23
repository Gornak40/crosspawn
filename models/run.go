package models

import "gorm.io/gorm"

type Run struct {
	gorm.Model

	EjudgeID        uint    `gorm:"unique;not null"`
	EjudgeContestID uint    `gorm:"not null"`
	EjudgeUserID    uint    `gorm:"not null"`
	EjudgeSource    string  `gorm:"not null"`
	EjudgeContest   Contest `gorm:"foreignKey:EjudgeContestID;not null"`

	ReviewCount uint `gorm:"not null"`
}
