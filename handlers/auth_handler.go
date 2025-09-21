package handlers

import (
	"CartoonBurgers/repositories"
	"CartoonBurgers/services"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (a *AuthHandlers) GetMiddlewareHadnler(ctx *gin.Context) gin.HandlerFunc {
	return services.AuthMiddleware(string(a.JwtKey))
}
