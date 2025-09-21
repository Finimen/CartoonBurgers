package handlers

import (
	"CartoonBurgers/repositories"
	"CartoonBurgers/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMenuHandler(ctx *gin.Context) {
	fmt.Println("MENU INITED")
	menuRepository, err := repositories.NewMenuRepository(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	menuService := services.NewMenuService(*menuRepository)

	products, err := menuService.GetMenu(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, products)
}
