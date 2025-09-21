package handlers

import (
	"CartoonBurgers/repositories"
	"CartoonBurgers/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandlers struct {
	Hasher *services.BcryptHasher
	JwtKey []byte
}

func NewAuthHandlers(jwtKey string) *AuthHandlers {
	return &AuthHandlers{
		Hasher: &services.BcryptHasher{},
		JwtKey: []byte(jwtKey),
	}
}

func (a *AuthHandlers) GetRegisterHandler(c *gin.Context) {
	repo, err := repositories.NewUserRepository(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Init Failed: " + err.Error()})
		return
	}

	registerHandler := services.NewRegisterHandler(a.Hasher, *repo)
	registerHandler.RegisterHandlerGin(c)
}

func (a *AuthHandlers) GetLoginHandler(c *gin.Context) {
	repo, err := repositories.NewUserRepository(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Init Failed: " + err.Error()})
		return
	}

	loginHandler := services.NewLoginHandler(a.Hasher, *repo, a.JwtKey)
	loginHandler.LoginHandlerGin(c)
}

func (a *AuthHandlers) AuthRequired() gin.HandlerFunc {
	return services.AuthMiddleware(string(a.JwtKey))
}

func (a *AuthHandlers) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr != "" {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			// Пытаемся распарсить токен
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
