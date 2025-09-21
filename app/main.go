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

	var auth = handlers.NewAuthHandlers("test_key")

	api := r.Group("/api")
	{
		api.GET("/menu", handlers.GetMenuHandler)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.GetRegisterHandler)
			authGroup.POST("/login", auth.GetLoginHandler)
		}

		//api.GET("/profile", auth.GetRegisterHandler)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":8080")
}
