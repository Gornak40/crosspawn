package controller

import (
	"errors"
	"fmt"
	"net/http"

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
		"Title": "Admin",
		"User":  user,
	})
}

func (s *Server) AdminPOST(c *gin.Context) {
	var form adminForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)
	user := session.Get("user")

	if err := s.validateJWT(form.JWT, user.(string)); err != nil { //nolint:forcetypeassert // it's ok
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

		return
	}

	session.Set("admin", true)
	_ = session.Save()

	c.Redirect(http.StatusFound, "/manage")
}

func (s *Server) validateJWT(t, user string) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(_ *jwt.Token) (interface{}, error) {
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

// TODO: check jwt here, it can expire.
func (s *Server) adminMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("admin")

	if admin == nil {
		c.Redirect(http.StatusFound, "/admin")
		c.Abort()

		return
	}

	c.Next()
}
