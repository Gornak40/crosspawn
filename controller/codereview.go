package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) CodereviewGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	submitEn, ok := session.Get("submit").(string)
	if !ok {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Select contest and problem",
			Type:    alerts.TypeInfo,
		})
		c.Redirect(http.StatusFound, "/")

		return
	}

	var submit reviewContext
	if err := json.Unmarshal([]byte(submitEn), &submit); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to unmarshal submit"})

		return
	}

	// TODO: get code from ejudge
	c.HTML(http.StatusOK, "codereview.html", gin.H{
		"Title":   "Review",
		"User":    user,
		"Submit":  submit,
		"Flashes": alerts.Get(session),
	})
}
