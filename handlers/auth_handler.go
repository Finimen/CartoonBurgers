package handlers

import (
	"CartoonBurgers/repositories"
	"CartoonBurgers/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandlers struct {
	Hasher   *services.BcryptHasher
	JwtKey   []byte
	userRepo *repositories.UserRepository
}

func NewAuthHandlers(jwtKey string, userRepo *repositories.UserRepository) *AuthHandlers {
	return &AuthHandlers{
		Hasher:   &services.BcryptHasher{},
		JwtKey:   []byte(jwtKey),
		userRepo: userRepo,
	}
}
func (a *AuthHandlers) GetRegisterHandler(c *gin.Context) {
	registerHandler := services.NewRegisterHandler(a.Hasher, *a.userRepo)
	registerHandler.RegisterHandlerGin(c)
}

func (a *AuthHandlers) GetLoginHandler(c *gin.Context) {
	loginHandler := services.NewLoginHandler(a.Hasher, *a.userRepo, a.JwtKey)
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
