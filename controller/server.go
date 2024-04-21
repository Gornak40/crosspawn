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
	db *gorm.DB
	ej *ejudge.EjClient
}

func NewServer(db *gorm.DB, ej *ejudge.EjClient) *Server {
	return &Server{
		db: db,
		ej: ej,
	}
}

func (s *Server) InitRouter(cfg *config.GinConfig) *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(cfg.Secret))
	r.Use(sessions.Sessions(sessionName, store))

	r.LoadHTMLGlob("./templates/*")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	r.GET("/", s.IndexGET)
	r.GET("/codereview", s.CodereviewGET)
	r.GET("/login", s.LoginGET)

	r.POST("/login", s.LoginPOST)
	r.POST("/logout", s.LogoutPOST)
	r.POST("/", s.IndexPOST)

	return r
}
