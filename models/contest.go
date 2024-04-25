package models

import (
	"strings"

	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model

	EjudgeID           uint   `gorm:"unique;not null"`
	EjudgeName         string `gorm:"not null;type:varchar(128)"`
	EjudgeProblemsList string `gorm:"not null;type:varchar(128)"`

	ReviewActive bool `gorm:"not null;default:true"`
	MaxRunID     uint `gorm:"not null"` // not inclusive
}

func NewContestFromEj(contest *ejudge.EjContest) *Contest {
	problemsList := make([]string, 0, len(contest.Problems))
	for _, p := range contest.Problems {
		problemsList = append(problemsList, p.ShortName)
	}

	return &Contest{
		EjudgeID:           contest.Contest.ID,
		EjudgeName:         contest.Contest.Name,
		EjudgeProblemsList: strings.Join(problemsList, " "),
	}
}
