package main

import (
	"github.com/Gornak40/crosspawn/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("./templates/*")

	r.StaticFile("/favicon.ico", "./static/img/favicon.ico")

	controller.Router(r)

	if err := r.Run(); err != nil {
		panic(err)
	}
}
