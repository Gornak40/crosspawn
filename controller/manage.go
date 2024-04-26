package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type manageForm struct {
	ContestID uint `binding:"required" form:"ejContestID"`
}

func (s *Server) ManageGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	var contests []models.Contest
	if err := s.db.Find(&contests).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.HTML(http.StatusOK, "manage.html", gin.H{
		"Title":    "Manage",
		"User":     user,
		"Contests": contests,
		"Flashes":  alerts.Get(session),
	})
}

// TODO: add check for judge credentials.
// TODO: add check for acm format.
func (s *Server) ManagePOST(c *gin.Context) {
	var form manageForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	contest, err := s.ej.GetContestStatus(form.ContestID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	dbContest := models.NewContestFromEj(contest)
	if err := s.db.Create(dbContest).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.Redirect(http.StatusFound, "/manage")
}

func (s *Server) ManageFlipPOST(c *gin.Context) {
	var form manageForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	dbContest := models.Contest{EjudgeID: form.ContestID}
	if err := s.db.Where(&dbContest).First(&dbContest).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	dbContest.ReviewActive = !dbContest.ReviewActive
	if err := s.db.Save(&dbContest).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.Redirect(http.StatusFound, "/manage")
}
