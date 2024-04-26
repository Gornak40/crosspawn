package models

import (
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model

	EjudgeID        uint   `gorm:"not null"`
	EjudgeContestID uint   `gorm:"not null"`
	EjudgeUserLogin string `gorm:"not null;type:varchar(128)"`
	EjudgeName      string `gorm:"not null;type:varchar(32)"`

	ReviewCount uint `gorm:"not null"`
}

func NewRunFromEj(run *ejudge.EjRun) *Run {
	return &Run{
		EjudgeID:        run.RunID,
		EjudgeContestID: run.ContestID,
		EjudgeUserLogin: run.UserLogin,
		EjudgeName:      run.ProbName,
	}
}
