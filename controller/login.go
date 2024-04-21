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

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login",
		"User":  session.Get("user"),
	})
}

// TODO: add auth.
func (s *Server) authUser(_, _ string) bool {
	return true
}

func (s *Server) LoginPOST(c *gin.Context) {
	session := sessions.Default(c)

	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.Redirect(http.StatusBadRequest, "/login")

		return
	}

	if !s.authUser(form.Login, form.Password) {
		c.Redirect(http.StatusUnauthorized, "/login")

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
