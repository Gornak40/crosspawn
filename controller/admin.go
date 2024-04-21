package controller

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	ErrForeignUser = errors.New("foreign user")
)

type adminForm struct {
	JWT string `binding:"required" form:"jwt"`
}

func (s *Server) AdminGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"Title": "Admin",
		"User":  user,
	})
}

func (s *Server) AdminPOST(c *gin.Context) {
	var form adminForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)
	user := session.Get("user")

	if err := parseUserJWT(form.JWT, user.(string)); err != nil { //nolint:forcetypeassert // it's ok
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	session.Set("admin", true)
	_ = session.Save()

	c.Redirect(http.StatusFound, "/manage")
}

// TODO: check JWT.
func parseUserJWT(_, _ string) error {
	return nil
}
