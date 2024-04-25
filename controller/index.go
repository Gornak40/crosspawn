package controller

import (
	"net/http"

	"github.com/Gornak40/crosspawn/internal/alerts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type reviewForm struct {
	Contest int    `binding:"required" form:"ejContest"`
	Problem string `binding:"required" form:"ejProblem"`
}

func (s *Server) IndexGET(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	contest := session.Get("contest")
	problem := session.Get("problem")

	if contest != nil && problem != nil {
		c.Redirect(http.StatusFound, "/codereview")

		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":   "Home",
		"User":    user,
		"Flashes": alerts.Get(session),
	})
}

func (s *Server) IndexPOST(c *gin.Context) {
	var form reviewForm
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	session := sessions.Default(c)

	session.Set("contest", form.Contest)
	session.Set("problem", form.Problem)
	_ = session.Save()
	c.Redirect(http.StatusFound, "/codereview")
}
