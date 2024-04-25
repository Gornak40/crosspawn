package models

import (
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model

	EjudgeID        uint    `gorm:"not null"`
	EjudgeContestID uint    `gorm:"not null"`
	EjudgeUserLogin string  `gorm:"not null;type:varchar(128)"`
	EjudgeSource    string  `gorm:"not null"`
	EjudgeContest   Contest `gorm:"foreignKey:EjudgeContestID;not null"`

	ReviewCount uint `gorm:"not null"`
}

func NewRunFromEj(run *ejudge.EjRun, source string) *Run {
	return &Run{
		EjudgeID:        run.RunID,
		EjudgeContestID: run.ContestID,
		EjudgeUserLogin: run.UserLogin,
		EjudgeSource:    source,
	}
}
