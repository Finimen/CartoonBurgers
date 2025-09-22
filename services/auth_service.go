package services

import (
	"CartoonBurgers/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IRedisClient interface {
	Get(key string) *redis.StringCmd
	Exists(key string) *redis.IntCmd
}

// @Summary User login implementation
// @Description Internal login handler service
func (l *LoginHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		l.Logger.Warn("invalid input format in login",
			"error", err.Error(),
			"client_ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if user.Username == "" || user.Password == "" {
		l.Logger.Warn("empty username or password in login",
			"client_ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	l.Logger.Debug("attempting login",
		"username", user.Username,
		"client_ip", c.ClientIP())

	storedPassword, err := l.Repository.GetUserByUsername(c.Request.Context(), user.Username)
	if err != nil {
		l.Logger.Warn("user not found or database error",
			"username", user.Username,
			"error", err.Error(),
			"client_ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credits"})
		return
	}

	err = l.Hasher.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		l.Logger.Warn("invalid password attempt",
			"username", user.Username,
			"client_ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credits"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(l.JwtKey)
	if err != nil {
		l.Logger.Error("failed to generate JWT token",
			"username", user.Username,
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generating error"})
		return
	}

	l.Logger.Info("user logged in successfully",
		"username", user.Username,
		"client_ip", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// @Summary User registration implementation
// @Description Internal registration handler service
func (h *RegisterHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.Logger.Warn("invalid input format in registration",
			"error", err.Error(),
			"client_ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "Invalid input"})
		return
	}

	if user.Username == "" || user.Password == "" || user.Email == "" {
		h.Logger.Warn("missing required fields in registration",
			"client_ip", c.ClientIP(),
			"has_username", user.Username != "",
			"has_password", user.Password != "",
			"has_email", user.Email != "")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	h.Logger.Debug("attempting user registration",
		"username", user.Username,
		"email", user.Email,
		"client_ip", c.ClientIP())

	hashedPassword, err := h.Hasher.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Logger.Error("password hashing failed",
			"username", user.Username,
			"error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hashing password error"})
		return
	}

	err = h.Repository.CreateUser(c.Request.Context(), user, string(hashedPassword))
	if err != nil {
		h.Logger.Warn("user registration failed - username already exists",
			"username", user.Username,
			"error", err.Error(),
			"client_ip", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exist"})
		return
	}

	h.Logger.Info("user registered successfully",
		"username", user.Username,
		"email", user.Email,
		"client_ip", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{"massage": "User registered successfully"})
}

// @Summary JWT authentication middleware
// @Description Middleware for validating JWT tokens and checking blacklist
func AuthMiddleware(jwtKey string, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")

		logger.Debug("auth middleware processing",
			"client_ip", c.ClientIP(),
			"has_authorization_header", tokenStr != "")

		if tokenStr == "" {
			logger.Warn("missing authorization token",
				"client_ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		redisClient := c.MustGet("redis").(IRedisClient)
		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := hex.EncodeToString(hash[:])
		exists, err := redisClient.Exists("blacklist:" + tokenHash).Result()

		if exists == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			c.Abort()
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signed method")
			}
			return []byte(jwtKey), nil
		})

		if err != nil {
			logger.Error("redis check failed",
				"error", err.Error(),
				"client_ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Warn("invalid token signature",
				"client_ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, exists := claims["username"]; exists {
				c.Set("username", username)
				logger.Debug("token validated successfully",
					"username", username.(string),
					"client_ip", c.ClientIP())
			} else {
				logger.Warn("username in token is not a string",
					"client_ip", c.ClientIP())
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
