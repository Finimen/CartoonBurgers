package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

type CartHandler struct {
	cookieSequre bool
}

func NewCartHandler(cookieSequre bool) *CartHandler {
	return &CartHandler{cookieSequre: cookieSequre}
}

func (h *CartHandler) getCartKey(c *gin.Context) string {
	token, exists := c.Get("token")
	if exists && token != "" {
		return fmt.Sprintf("cart:user:%s", token)
	}

	sessionID, err := c.Cookie("cart_session")
	if err != nil || sessionID == "" {
		sessionID = generateSessionID()
		c.SetCookie("cart_session", sessionID, 30*24*3600, "/", "", h.cookieSequre, true)
	}
	return fmt.Sprintf("cart:session:%s", sessionID)
}

// @Summary Get all products
// @Description For menu from db
// @Tags cart
// @Produce json
// @Param Context
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Token missing"
// @Router [get]
func (h *CartHandler) GetCartHandler(c *gin.Context) {
	redisClient := c.MustGet("redis").(*redis.Client)
	cartKey := h.getCartKey(c)

	cartData, err := redisClient.Get(cartKey).Result()

	if err == redis.Nil {
		c.JSON(200, gin.H{"items": []CartItem{}})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Getting cart error"})
		return
	}

	var cart []CartItem
	json.Unmarshal([]byte(cartData), &cart)
	c.JSON(200, gin.H{"items": cart})
}

// @Summary Get all products
// @Description For menu from db
// @Tags cart
// @Produce json
// @Param Context
// @Success 200 {object} []CartItem "Item added to cart"
// @Failure 400 {object} gin.H "Invalid format"
// @Failure 400 {object} gin.H "Cart error"
// @Failure 400 {object} gin.H "Saving Cart error"
// @Router /add [post]
func (h *CartHandler) AddToCartHandler(c *gin.Context) {
	redisClient := c.MustGet("redis").(*redis.Client)
	var item CartItem

	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	cartKey := h.getCartKey(c)

	cartData, err := redisClient.Get(cartKey).Result()
	var cart []CartItem

	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cart error"})
		return
	}

	if err != redis.Nil {
		json.Unmarshal([]byte(cartData), &cart)
	}

	found := false
	for i, cartItem := range cart {
		if cartItem.ProductID == item.ProductID {
			cart[i].Quantity += item.Quantity
			found = true
			break
		}
	}

	if !found {
		cart = append(cart, item)
	}

	cartJSON, _ := json.Marshal(cart)
	err = redisClient.Set(cartKey, cartJSON, 30*24*time.Hour).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Saving Cart error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в корзину", "cart": cart})
}

// @Summary Remove from cart
// @Tags cart
// @Produce json
// @Param Context
// @Success 200 {object} []CartItem "Item removed from cart"
// @Failure 400 {object} gin.H "Invalid format"
// @Failure 400 {object} gin.H "Cart error"
// @Failure 400 {object} gin.H "Removing Cart error"
// @Router /:productId [delete]
func (h *CartHandler) RemoveFromCartHandler(c *gin.Context) {
	redisClient := c.MustGet("redis").(*redis.Client)
	var item CartItem

	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	cartKey := h.getCartKey(c)

	cartData, err := redisClient.Get(cartKey).Result()
	var cart []CartItem

	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cart error"})
		return
	}

	if err != redis.Nil {
		json.Unmarshal([]byte(cartData), &cart)
	}

	for i, cartItem := range cart {
		if cartItem.ProductID == item.ProductID {
			cart = append(cart[:i], cart[i+1:]...)
			break
		}
	}

	cartJSON, _ := json.Marshal(cart)
	err = redisClient.Set(cartKey, cartJSON, 30*24*time.Hour).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Removing Cart error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар убран из корзины", "cart": cart})
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
