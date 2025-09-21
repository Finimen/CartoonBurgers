package handlers

import (
	"CartoonBurgers/repositories"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfileHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Создаем репозиторий для работы с БД
	repo, err := repositories.NewUserRepository(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection error"})
		return
	}

	// Получаем данные пользователя из БД
	user, err := repo.GetUserProfile(c.Request.Context(), username.(string))
	if err != nil {
		fmt.Println("ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ", username.(string))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"bonuses":  user.Bonus,
	})
}
