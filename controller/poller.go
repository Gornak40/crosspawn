package controller

import (
	"github.com/Gornak40/crosspawn/models"
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"github.com/sirupsen/logrus"
)

func (s *Server) Poll() error {
	var contests []models.Contest
	if err := s.db.Where(&models.Contest{ReviewActive: true}).Find(&contests).Error; err != nil {
		return err
	}

	for _, contest := range contests {
		logrus.WithField("contestID", contest.EjudgeID).Info("polling contest")
		if err := s.pollContest(&contest); err != nil { //nolint:gosec // G601: Implicit memory aliasing in for loop.
			logrus.WithError(err).WithField("contestID", contest.EjudgeID).Errorf("failed to poll contest")
		}
	}

	return nil
}

func statusFromRating(rating int) ejudge.RunStatus {
	if rating > 0 { // TODO: think about it
		return ejudge.RunStatusOK
	}

	return ejudge.RunStatusRJ
}

func runStatusFromDB(dbStatus *models.Run) *ejudge.EjStatusChange {
	return &ejudge.EjStatusChange{
		RunID:     dbStatus.EjudgeID,
		ContestID: dbStatus.EjudgeContestID,
		Status:    statusFromRating(dbStatus.Rating),
	}
}

func (s *Server) pollContest(dbContest *models.Contest) error {
	runs, err := s.ej.GetContestRuns(dbContest.EjudgeID, "status == PR")
	if err != nil {
		return err
	}
	logrus.WithField("count", len(runs.Runs)).Info("ejudge runs received")
	runsMapa := make(map[uint]struct{})
	for _, run := range runs.Runs {
		runsMapa[run.RunID] = struct{}{}
	}

	var dbRuns []models.Run
	quRun := models.Run{EjudgeContestID: dbContest.EjudgeID}
	if err := s.db.Where(&quRun).Find(&dbRuns).Error; err != nil {
		return err
	}

	// change status of runs
	oldMapa := make(map[uint]struct{})
	for _, dbRun := range dbRuns {
		if _, ok := runsMapa[dbRun.EjudgeID]; !ok {
			logrus.WithField("runID", dbRun.EjudgeID).Warn("run is lost in ejudge")
			if err := s.db.Delete(&dbRun).Error; err != nil { //nolint:gosec // G601: Implicit memory aliasing in for loop.
				logrus.WithError(err).Error("failed to delete run")
			}

			continue
		}
		if dbRun.ReviewCount >= uint(s.cfg.ReviewLimit) {
			status := runStatusFromDB(&dbRun) //nolint:gosec // G601: Implicit memory aliasing in for loop.
			if err := s.ej.ChangeRunStatus(status); err != nil {
				logrus.WithError(err).Error("failed to change run status")
			}
			logrus.WithField("runID", dbRun.EjudgeID).Info("run review is done")
			if err := s.db.Delete(&dbRun).Error; err != nil { //nolint:gosec // G601: Implicit memory aliasing in for loop.
				logrus.WithError(err).Error("failed to delete run")
			}
		}
		oldMapa[dbRun.EjudgeID] = struct{}{}
	}

	// save new runs
	// TODO: fix saving deleted runs
	for _, run := range runs.Runs {
		if _, ok := oldMapa[run.RunID]; ok {
			continue
		}
		dbRun := models.NewRunFromEj(&run) //nolint:gosec // G601: Implicit memory aliasing in for loop.
		logrus.WithField("run", run).Info("saving new run")
		if err := s.db.Save(&dbRun).Error; err != nil {
			logrus.WithError(err).Error("failed to save run")
		}
	}

	return nil
}
