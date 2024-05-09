package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type codereviewFrom struct {
	RatingDelta int `binding:"required" form:"ratingDelta"`
}

func (s *Server) CodereviewGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	submit, _ := session.Get("submit").(reviewContext)

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
	submit, _ := session.Get("submit").(reviewContext)
	user, _ := session.Get("user").(string)

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

	var dbUser models.User
	quUser := models.User{EjudgeLogin: user}
	if err := s.db.Where(&quUser).First(&dbUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	switch {
	case form.RatingDelta > 0:
		dbUser.ReviewAproveCount++
	case form.RatingDelta < 0:
		dbUser.ReviewRejectCount++
	}

	if err := s.db.Save(&dbUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

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
