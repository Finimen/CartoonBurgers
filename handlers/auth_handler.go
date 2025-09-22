package handlers

import (
	"CartoonBurgers/repositories"
	"CartoonBurgers/services"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandlers struct {
	Hasher   *services.BcryptHasher
	JwtKey   []byte
	userRepo *repositories.UserRepository
	logger   *slog.Logger
}

func NewAuthHandlers(jwtKey string, userRepo *repositories.UserRepository, logger *slog.Logger) *AuthHandlers {
	return &AuthHandlers{
		Hasher:   &services.BcryptHasher{},
		JwtKey:   []byte(jwtKey),
		userRepo: userRepo,
		logger:   logger,
	}
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User registration data"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (a *AuthHandlers) GetRegisterHandler(c *gin.Context) {
	registerHandler := services.NewRegisterHandler(a.Hasher, *a.userRepo, a.logger)
	registerHandler.RegisterHandlerGin(c)
}

// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "Login credentials"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (a *AuthHandlers) GetLoginHandler(c *gin.Context) {
	loginHandler := services.NewLoginHandler(a.Hasher, *a.userRepo, a.JwtKey, a.logger)
	loginHandler.LoginHandlerGin(c)
}

// @Summary User logout
// @Description Invalidate user's JWT token
// @Tags auth
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Token missing"
// @Failure 500 {object} map[string]string "Logout failed"
// @Router /logout [post]
func (a *AuthHandlers) GetLogoutHandler(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token missing"})
		return
	}

	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])

	redisClient := c.MustGet("redis").(*redis.Client)
	token, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.JwtKey), nil
	})

	var expiration time.Duration
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := time.Unix(int64(claims["exp"].(float64)), 0)
		expiration = time.Until(exp)
	} else {
		expiration = time.Hour
	}

	err := redisClient.SetNX("blacklist:"+tokenHash, "1", expiration).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// @Summary Authentication required middleware
// @Description JWT authentication middleware that requires valid token
// @Tags auth
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
func (a *AuthHandlers) AuthRequired() gin.HandlerFunc {
	return services.AuthMiddleware(string(a.JwtKey), a.logger)
}

// @Summary Optional authentication middleware
// @Description JWT authentication middleware that works with or without token
// @Tags auth
// @Security ApiKeyAuth
// @Param Authorization header string false "Bearer token (optional)"
func (a *AuthHandlers) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return a.JwtKey, nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					c.Set("username", claims["username"])
					c.Set("token", tokenStr)
				}
			}
		}
		c.Next()
	}
}
