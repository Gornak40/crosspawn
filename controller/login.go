package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login",
	})
}
