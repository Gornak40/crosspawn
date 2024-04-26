package controller

import (
	"encoding/gob"

	"github.com/Gornak40/crosspawn/config"
	"github.com/Gornak40/crosspawn/pkg/ejudge"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const sessionName = "crosspawn"

type Server struct {
	db  *gorm.DB
	ej  *ejudge.EjClient
	cfg *config.ServerConfig
}

func NewServer(db *gorm.DB, ej *ejudge.EjClient, cfg *config.ServerConfig) *Server {
	return &Server{
		db:  db,
		ej:  ej,
		cfg: cfg,
	}
}

func (s *Server) InitRouter() *gin.Engine {
	r := gin.Default()

	gob.Register(reviewContext{})

	store := cookie.NewStore([]byte(s.cfg.GinSecret))
	r.Use(sessions.Sessions(sessionName, store))

	r.LoadHTMLGlob("./templates/*")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	{
		r.GET("/login", s.LoginGET)
		r.POST("/login", s.LoginPOST)
	}

	ua := r.Group("/", s.userMiddleware)
	{
		ua.GET("/", s.IndexGET)
		ua.POST("/", s.IndexPOST)

		ua.GET("/admin", s.AdminGET)
		ua.POST("/admin", s.AdminPOST)

		ua.GET("/profile", s.ProfileGET)

		ua.POST("/logout", s.LogoutPOST)
	}

	ca := ua.Group("/codereview", s.codereviewMiddleware)
	{
		ca.GET("/", s.CodereviewGET)
		ca.POST("/", s.CodereviewPOST)
	}

	aa := ua.Group("/manage", s.adminMiddleware)
	{
		aa.GET("/", s.ManageGET)
		aa.POST("/", s.ManagePOST)

		aa.POST("/flip", s.ManageFlipPOST)
	}

	return r
}
