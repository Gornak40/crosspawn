package controller

import (
	"errors"
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	errInvalidSessionSubmit = errors.New("invalid submit in session")
)

type codereviewFrom struct {
	RatingDelta int `binding:"required" form:"ratingDelta"`
}

func (s *Server) CodereviewGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	submit, ok := session.Get("submit").(reviewContext)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errInvalidSessionSubmit.Error()})

		return
	}

	source, err := s.ej.GetRunSource(submit.ContestID, submit.RunID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	submit.Source = source
	c.HTML(http.StatusOK, "codereview.html", gin.H{
		"Title":   "Review",
		"User":    user,
		"Submit":  submit,
		"Flashes": alerts.Get(session),
	})
}

func (s *Server) CodereviewPOST(c *gin.Context) {
	var form codereviewFrom

	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)

	submit, ok := session.Get("submit").(reviewContext)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errInvalidSessionSubmit.Error()})

		return
	}

	var dbRun models.Run
	quRun := models.Run{EjudgeID: submit.RunID, EjudgeContestID: submit.ContestID}
	if err := s.db.Where(&quRun).First(&dbRun).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	dbRun.ReviewCount++
	dbRun.Rating += form.RatingDelta

	if err := s.db.Save(&dbRun).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	session.Delete("submit")
	_ = session.Save()

	_ = alerts.Add(session, alerts.Alert{
		Message: "Review processed",
		Type:    alerts.TypeSuccess,
	})
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) codereviewMiddleware(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("submit") == nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Select contest and problem",
			Type:    alerts.TypeInfo,
		})
		c.Redirect(http.StatusFound, "/")
		c.Abort()

		return
	}

	c.Next()
}
