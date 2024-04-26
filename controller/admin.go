package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		"Title":   "Admin",
		"User":    user,
		"Flashes": alerts.Get(session),
	})
}

func (s *Server) AdminPOST(c *gin.Context) {
	var form adminForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)
	user, ok := session.Get("user").(string)

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})

		return
	}

	if err := s.validateJWT(form.JWT, user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	session.Set("jwt", form.JWT)
	_ = session.Save()

	_ = alerts.Add(session, alerts.Alert{ // TODO: add expiration time
		Message: "JWT is valid",
		Type:    alerts.TypeSuccess,
	})
	c.Redirect(http.StatusFound, "/manage")
}

func (s *Server) validateJWT(token, user string) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return err
	}
	if claims["sub"] != user {
		return fmt.Errorf("%w: token is not owned by %s", ErrForeignUser, user)
	}

	return nil
}

func (s *Server) adminMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	user, ok := session.Get("user").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})

		return
	}

	token, ok := session.Get("jwt").(string)
	if !ok {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Enter admin JWT",
			Type:    alerts.TypeWarning,
		})
		c.Redirect(http.StatusFound, "/admin")
		c.Abort()

		return
	}

	if err := s.validateJWT(token, user); err != nil {
		_ = alerts.Add(session, alerts.Alert{
			Message: "Your JWT is expired",
			Type:    alerts.TypeDanger,
		})
		c.Redirect(http.StatusFound, "/admin")
		c.Abort()

		return
	}

	c.Next()
}
