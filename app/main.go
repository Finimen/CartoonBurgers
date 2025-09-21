package main

import (
	"CartoonBurgers/handlers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func initRedis() *gin.Engine {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // или твой Redis URL
		Password: "",               // пароль если есть
		DB:       0,                // база данных
	})

	if err := rdb.Ping().Err(); err != nil {
		panic("Не удалось подключиться к Redis: " + err.Error())
	}

	r := gin.Default()

	// Добавляем Redis в контекст
	r.Use(func(c *gin.Context) {
		c.Set("redis", rdb)
		c.Next()
	})

	return r
}

func main() {
	limiter := NewRateLimiter(100, time.Minute)

	r := initRedis()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	var auth = handlers.NewAuthHandlers("test_key")

	api := r.Group("/api")
	api.Use(RateLimitMiddleware(limiter))
	{
		api.GET("/menu", handlers.GetMenuHandler)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.GetRegisterHandler)
			authGroup.POST("/login", auth.GetLoginHandler)
		}

		cartGroup := api.Group("/cart")
		cartGroup.Use(auth.OptionalAuth())
		{
			cartGroup.GET("", handlers.GetCartHandler)
			cartGroup.POST("/add", handlers.AddToCartHandler)
			cartGroup.DELETE("/:productId", handlers.RemoveFromCartHandler)
		}

		protected := api.Group("")
		protected.Use(auth.AuthRequired())
		{
			protected.GET("/profile", handlers.GetProfileHandler)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run(":8080")
}
