package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginForm struct {
	Login     string `binding:"required" form:"ejLogin"`
	Password  string `binding:"required" form:"ejPassword"`
	ContestID uint   `binding:"required" form:"ejContest"`
}

func (s *Server) LoginGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title":   "Login",
		"User":    user,
		"Flashes": alerts.Get(session),
	})
}

func (s *Server) LoginPOST(c *gin.Context) {
	session := sessions.Default(c)

	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := s.ej.AuthUser(ejudge.AuthHeader{
		Login:     form.Login,
		Password:  form.Password,
		ContestID: form.ContestID,
	}); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	session.Set("user", form.Login)
	_ = session.Save()

	_ = alerts.Add(session, alerts.Alert{
		Message: "You are logged in",
		Type:    alerts.TypeSuccess,
	})
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) LogoutPOST(c *gin.Context) {
	session := sessions.Default(c)

	session.Clear()
	_ = session.Save()
	c.Redirect(http.StatusFound, "/")
}

func (s *Server) userMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "You are not logged in",
			Type:    alerts.TypeWarning,
		})
		c.Redirect(http.StatusFound, "/login")
		c.Abort()

		return
	}

	c.Next()
}
