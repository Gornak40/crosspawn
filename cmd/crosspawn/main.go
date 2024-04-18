package main

import (
	"log"

	"github.com/Gornak40/crosspawn/controller"
	"github.com/Gornak40/crosspawn/models"
	"github.com/gin-gonic/gin"
)

func initRouter() (*gin.Engine, error) {
	db, err := models.ConnectDatabase()
	if err != nil {
		return nil, err
	}

	r := gin.Default()

	r.LoadHTMLGlob("./templates/*")
	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	s := controller.NewServer(db)

	r.GET("/", s.Index)
	r.GET("/codereview", s.Codereview)
	r.GET("/login", s.Login)

	return r, nil
}

func main() {
	r, err := initRouter()

	if err != nil {
		log.Fatal(err)
	}

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
