package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Login    string `binding:"required" form:"ejLogin"`
	Password string `binding:"required" form:"ejPassword"`
}

func (s *Server) LoginGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login",
		"User":  user,
	})
}

func (s *Server) LoginPOST(c *gin.Context) {
	session := sessions.Default(c)

	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if !s.authUser(form.Login, form.Password) {
		c.Redirect(http.StatusFound, "/login")

		return
	}

	session.Set("user", form.Login)
	_ = session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) LogoutPOST(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	_ = session.Save()
	c.Redirect(http.StatusFound, "/")
}

// TODO: add auth.
func (s *Server) authUser(_, _ string) bool {
	return true
}

func (s *Server) userMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()

		return
	}

	c.Next()
}
