package controller

import (
	"net/http"

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
	if res := s.db.Find(&contests); res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})

		return
	}

	c.HTML(http.StatusOK, "manage.html", gin.H{
		"Title":    "Manage",
		"User":     user,
		"Contests": contests,
	})
}

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
	if res := s.db.Create(dbContest); res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})

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
	if err := s.db.Where("ejudge_id = ?", form.ContestID).First(&dbContest).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	dbContest.ReviewActive = !dbContest.ReviewActive
	if res := s.db.Save(&dbContest); res.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})

		return
	}

	c.Redirect(http.StatusFound, "/manage")
}
