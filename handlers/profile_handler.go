package handlers

import (
	"CartoonBurgers/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userRepo *repositories.UserRepository
}

func NewProfileHandler(repo *repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: repo}
}

// @Summary Get profile info from user
// @Tags profile
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} gin.H "User not authenticated"
// @Failure 500 {object} gin.H "Failed to get user profile"
// @Router /profile [get]
func (h *ProfileHandler) GetProfileHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	repo := h.userRepo

	user, err := repo.GetUserProfile(c.Request.Context(), username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
		"bonuses":  user.Bonus,
	})
}
