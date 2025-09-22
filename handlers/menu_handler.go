package handlers

import (
	"CartoonBurgers/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repositories.MenuRerository
}

func NewMenuHandler(repo *repositories.MenuRerository) *MenuHandler {
	return &MenuHandler{menuRepo: repo}
}

func (h *MenuHandler) GetMenu(ctx *gin.Context) {
	products, err := h.menuRepo.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}
