package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) ProfileGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"Title":   "Profile GET",
		"User":    user,
		"Flashes": alerts.Get(session),
	})
}
