package main

import (
	"CartoonBurgers/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	api := r.Group("/api")
	{
		api.GET("/menu", handlers.GetMenuHandler)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":8080")
}
