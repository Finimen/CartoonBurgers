package handlers

import (
	"CartoonBurgers/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repositories.ProductRerository
}

func NewMenuHandler(repo *repositories.ProductRerository) *MenuHandler {
	return &MenuHandler{menuRepo: repo}
}

// @Summary Get restaurant menu
// @Tags menu
// @Produce json
// @Success 200 {object} []models.Products
// @Router /menu [get]
func (h *MenuHandler) GetMenu(ctx *gin.Context) {
	products, err := h.menuRepo.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}
