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
		logrus.WithField("contestID", contest.EjudgeID).Info("polling contest")
		if err := s.pollContest(&contest); err != nil { //nolint:gosec // G601: Implicit memory aliasing in for loop.
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
		logrus.Info(run)                              // TODO: remove
		dbRun := models.NewRunFromEj(&run, "babayka") //nolint:gosec // G601: Implicit memory aliasing in for loop.
		if res := s.db.Create(dbRun); res.Error != nil {
			logrus.WithError(err).WithFields(logrus.Fields{"contestID": run.ContestID, "runID": run.RunID}).
				Error("failed to save run")
		}
	}
	dbContest.MaxRunID = min(to, runs.TotalRuns)

	return s.db.Save(dbContest).Error
}
