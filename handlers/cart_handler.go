package handlers

import (
	"encoding/json"
	"fmt"
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

func (h *CartHandler) AddToCartHandler(c *gin.Context) {
	redisClient := c.MustGet("redis").(*redis.Client)
	var item CartItem

	if err := c.BindJSON(&item); err != nil {
		c.JSON(400, gin.H{"error": "Invalid format"})
		return
	}

	cartKey := h.getCartKey(c)

	cartData, err := redisClient.Get(cartKey).Result()
	var cart []CartItem

	if err != nil && err != redis.Nil {
		c.JSON(500, gin.H{"error": "Cart error"})
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
		c.JSON(500, gin.H{"error": "Saving Cart error"})
		return
	}

	c.JSON(200, gin.H{"message": "Товар добавлен в корзину", "cart": cart})
}

func (h *CartHandler) RemoveFromCartHandler(c *gin.Context) {

}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
