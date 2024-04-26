package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	errNoOKSubmit = errors.New("no OK submit")
)

type reviewForm struct {
	Contest uint   `binding:"required" form:"ejContest"`
	Problem string `binding:"required" form:"ejProblem"`
}

func (s *Server) IndexGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if session.Get("submit") != nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Finish your current review",
			Type:    alerts.TypeWarning,
		})
		c.Redirect(http.StatusFound, "/codereview")

		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":   "Home",
		"User":    user,
		"Flashes": alerts.Get(session),
	})
}

type reviewContext struct {
	RunID     uint   `json:"run_id"`
	ContestID uint   `json:"contest_id"`
	Problem   string `json:"problem"`
	Source    string `json:"source"`
}

//nolint:funlen // Life is hard
func (s *Server) IndexPOST(c *gin.Context) {
	var form reviewForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)
	user, ok := session.Get("user").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})

		return
	}

	// User must not have current review
	if session.Get("submit") != nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "user has unclosed review"})

		return
	}

	// Contest must be valid
	if s.db.Where(&models.Contest{EjudgeID: form.Contest, ReviewActive: true}).First(&models.Contest{}).Error != nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: fmt.Sprintf("Contest %d is not in pool", form.Contest),
			Type:    alerts.TypeDanger,
		})
		c.Redirect(http.StatusFound, "/")

		return
	}

	// User must have OK submit
	if err := s.reviewValid(user, form.Contest, form.Problem); err != nil {
		if errors.Is(err, errNoOKSubmit) {
			_ = alerts.Add(session, alerts.Alert{
				Message: fmt.Sprintf("Solve %s in %d first", form.Problem, form.Contest),
				Type:    alerts.TypeDanger,
			})
			c.Redirect(http.StatusFound, "/")

			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// Good submit must exist
	var dbRun models.Run
	magicFilter := "ejudge_contest_id = ? AND ejudge_name = ? AND ejudge_user_login != ? AND review_count <= ?"
	if err := s.db.Where(magicFilter, form.Contest, form.Problem, user, s.cfg.ReviewLimit).
		Order("RANDOM()").First(&dbRun).Error; err != nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "No runs to review",
			Type:    alerts.TypeWarning,
		})
		c.Redirect(http.StatusFound, "/")

		return
	}

	rctx := reviewContext{
		RunID:     dbRun.EjudgeID,
		ContestID: dbRun.EjudgeContestID,
		Problem:   dbRun.EjudgeName,
	}
	data, _ := json.Marshal(rctx) //nolint:errchkjson // It's my JSON
	session.Set("submit", string(data))
	_ = session.Save()

	_ = alerts.Add(session, alerts.Alert{
		Message: "Review started",
		Type:    alerts.TypeSuccess,
	})
	c.Redirect(http.StatusFound, "/codereview")
}

func (s *Server) reviewValid(user string, contestID uint, problem string) error {
	filter := fmt.Sprintf("login == '%s' && prob == '%s' && (status == OK || status == PR)", user, problem)
	runs, err := s.ej.GetContestRuns(contestID, filter)
	if err != nil {
		return err
	}
	if len(runs.Runs) == 0 {
		return fmt.Errorf("%w: %d", errNoOKSubmit, contestID)
	}

	return nil
}
