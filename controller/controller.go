package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/", index)
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
