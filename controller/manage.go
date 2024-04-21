package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) ManageGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	admin := session.Get("admin")

	if admin == nil { // TODO: move this to middleware
		c.Redirect(http.StatusFound, "/admin")

		return
	}

	c.HTML(http.StatusOK, "manage.html", gin.H{
		"Title": "Manage",
		"User":  user,
	})
}
