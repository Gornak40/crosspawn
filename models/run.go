package models

import (
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model

	EjudgeID        uint   `gorm:"not null;uniqueIndex:idx_full_id"`
	EjudgeContestID uint   `gorm:"not null;uniqueIndex:idx_full_id"`
	EjudgeUserLogin string `gorm:"not null;type:varchar(128)"`
	EjudgeName      string `gorm:"not null;type:varchar(32)"`

	ReviewCount uint `gorm:"not null"`
	Rating      int  `gorm:"not null"`
}

func NewRunFromEj(run *ejudge.EjRun) *Run {
	return &Run{
		EjudgeID:        run.RunID,
		EjudgeContestID: run.ContestID,
		EjudgeUserLogin: run.UserLogin,
		EjudgeName:      run.ProbName,
	}
}
