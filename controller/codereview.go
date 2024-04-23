package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) CodereviewGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	contest := session.Get("contest")
	problem := session.Get("problem")

	if contest == nil || problem == nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	// TODO: get code from ejudge
	c.HTML(http.StatusOK, "codereview.html", gin.H{
		"Title":     "Review",
		"User":      user,
		"CodeTitle": contest,
		"Code":      problem,
	})
}
