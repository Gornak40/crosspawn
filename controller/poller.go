package controller

import (
	"fmt"

	"github.com/Gornak40/crosspawn/models"
	"github.com/sirupsen/logrus"
)

func (s *Server) Poll() error {
	var contests []models.Contest
	if res := s.db.Where("review_active = 1").Find(&contests); res.Error != nil {
		return res.Error
	}

	for _, contest := range contests {
		contest := contest
		logrus.WithField("contestID", contest.EjudgeID).Info("polling contest")
		if err := s.pollContest(&contest); err != nil {
			logrus.WithError(err).WithField("contestID", contest.EjudgeID).Errorf("failed to poll contest")
		}
	}

	return nil
}

func (s *Server) pollContest(dbContest *models.Contest) error {
	from := dbContest.MaxRunID
	to := from + uint(s.cfg.PollBatchSize)
	filter := fmt.Sprintf("status == PR && %d <= id && id < %d", from, to)

	runs, err := s.ej.GetContestRuns(dbContest.EjudgeID, filter)
	if err != nil {
		return err
	}

	for _, run := range runs.Runs {
		run := run
		logrus.Info(run)
		dbRun := models.NewRunFromEj(&run, "babayka")
		if res := s.db.Create(dbRun); res.Error != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"contestID": run.ContestID, "runID": run.RunID}).
				Error("failed to save run")
		}
	}
	dbContest.MaxRunID = min(to, runs.TotalRuns)

	return s.db.Save(dbContest).Error
}
