package main

import (
	"CartoonBurgers/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	r.GET("/api/menu")

	api := r.Group("/api")
	{
		api.GET("/menu", handlers.GetMenuHandler())
		api.GET("")
	}
}
