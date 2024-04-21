package controller

import (
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

	store := cookie.NewStore([]byte(s.cfg.GinSecret))
	r.Use(sessions.Sessions(sessionName, store))

	r.LoadHTMLGlob("./templates/*")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	r.GET("/login", s.LoginGET)
	r.POST("/login", s.LoginPOST)

	ua := r.Group("/")
	ua.Use(s.userMiddleware)
	{
		ua.GET("/", s.IndexGET)
		ua.GET("/codereview", s.CodereviewGET)
		ua.GET("/admin", s.AdminGET)
		ua.GET("/manage", s.ManageGET)

		ua.POST("/logout", s.LogoutPOST)
		ua.POST("/", s.IndexPOST)
		ua.POST("/admin", s.AdminPOST)
	}

	return r
}
