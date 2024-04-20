package controller

import (
	"net/http"

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

	if user == nil {
		c.Redirect(http.StatusFound, "/login")

		return
	}

	contest := session.Get("contest")
	problem := session.Get("problem")

	if contest != nil && problem != nil {
		c.Redirect(http.StatusFound, "/codereview")

		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Home",
		"User":  user,
	})
}

func (s *Server) IndexPOST(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user == nil {
		c.Redirect(http.StatusFound, "/login")

		return
	}

	var form reviewForm
	if err := c.ShouldBind(&form); err != nil {
		c.Redirect(http.StatusFound, "/")

		return
	}

	session.Set("contest", form.Contest)
	session.Set("problem", form.Problem)
	_ = session.Save()
	c.Redirect(http.StatusFound, "/codereview")
}
