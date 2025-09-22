package main

import (
	"CartoonBurgers/app/config"
	"CartoonBurgers/handlers"
	"CartoonBurgers/repositories"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	_ "CartoonBurgers/docs"
)

func initRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Fatal("Не удалось подключиться к Redis: ", err)
	}

	return rdb
}

// @title User API
// @version 1.0
// @description API for Cartoon Burgers authentication service
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @contact.name FinimenSniperC
// @contact.email finimensniper@gmail.com
// @BasePath /api/v1
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	var logger *slog.Logger
	if cfg.Environment.Current == "development" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)

	logger.Info("application starting",
		"environment", cfg.Environment.Current,
		"version", "1.0")

	appRepo, err := repositories.NewAppRepository(context.Background(), cfg.Database.Path)
	if err != nil {
		log.Fatal("Cannot init repositories:", err)
	}
	defer appRepo.DB.Close()

	rdb := initRedis(cfg.Redis)
	defer rdb.Close()

	limiter := NewRateLimiter(cfg.RateLimit.MaxRequests, cfg.RateLimit.Window)

	r := gin.Default()

	rAdapter := NewRedisAdapter(rdb)

	r.Use(func(c *gin.Context) {
		c.Set("redis", rAdapter)
		c.Next()
	})

	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	menuHandler := handlers.NewMenuHandler(appRepo.ProductRerository)
	authHandler := handlers.NewAuthHandlers(cfg.JWT.SecretKey, appRepo.UserRepository, logger)
	profileHandler := handlers.NewProfileHandler(appRepo.UserRepository)
	cartHandler := handlers.NewCartHandler(cfg.Server.CookieSecure)

	api := r.Group("/api")
	api.Use(RateLimitMiddleware(limiter))
	{
		api.GET("/menu", menuHandler.GetMenu)

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.GetRegisterHandler)
			authGroup.POST("/login", authHandler.GetLoginHandler)
			authGroup.POST("/logout", authHandler.GetLogoutHandler)
		}

		cartGroup := api.Group("/cart")
		cartGroup.Use(authHandler.OptionalAuth())
		{
			cartGroup.GET("", cartHandler.GetCartHandler)
			cartGroup.POST("/add", cartHandler.AddToCartHandler)
			cartGroup.DELETE("/:productId", cartHandler.RemoveFromCartHandler)
		}

		protected := api.Group("")
		protected.Use(authHandler.AuthRequired())
		{
			protected.GET("/profile", profileHandler.GetProfileHandler)
		}
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server is starting on %s\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
