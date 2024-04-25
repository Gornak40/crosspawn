package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) CodereviewGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	contest := session.Get("contest")
	problem := session.Get("problem")

	if contest == nil || problem == nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Please select contest and problem",
			Type:    alerts.TypeInfo,
		})
		c.Redirect(http.StatusFound, "/")

		return
	}

	// TODO: get code from ejudge
	c.HTML(http.StatusOK, "codereview.html", gin.H{
		"Title":     "Review",
		"User":      user,
		"CodeTitle": contest,
		"Code":      problem,
		"Flashes":   alerts.Get(session),
	})
}
