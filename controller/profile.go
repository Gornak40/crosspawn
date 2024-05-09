package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) ProfileGET(c *gin.Context) {
	session := sessions.Default(c)
	user, _ := session.Get("user").(string)

	var dbUser models.User
	quUser := models.User{EjudgeLogin: user}
	if err := s.db.Where(&quUser).First(&dbUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Title":           "Profile GET",
		"User":            user,
		"UserAproveCount": dbUser.ReviewAproveCount,
		"UserRejectCount": dbUser.ReviewRejectCount,
		"Flashes":         alerts.Get(session),
	})
}
