package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//RESPONSE STRUCTS

type CartItemResponse struct {
	ID        uint           `json:"id"`
	Product   ProductSummary `json:"product"`
	Quantity  int            `json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ProductSummary struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	ImageURL      string  `json:"image_url"`
}

//MAPPER
func mapCartItem(item models.CartItem) CartItemResponse {
	return CartItemResponse{
		ID:        item.ID,
		Quantity:  item.Quantity,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		Product: ProductSummary{
			ID:            item.Product.ID,
			Name:          item.Product.Name,
			Description:   item.Product.Description,
			Price:         item.Product.Price,
			StockQuantity: int(item.Product.StockQuantity),
			ImageURL:      item.Product.ImageURL,
		},
	}
}


// ADD TO CART 


func AddToCart(c *gin.Context) {

	userIDVal, exists := c.Get("userID")
	userID, ok := userIDVal.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if input.Quantity > int(product.StockQuantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock available"})
		return
	}

	var existing models.CartItem
	if err := config.DB.
		Where("user_id = ? AND product_id = ?", userID, input.ProductID).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, gin.H{"error": "product already in cart"})
		return
	}

	cartItem := models.CartItem{
		UserID:    userID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.
		Preload("Product").
		First(&cartItem, cartItem.ID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cart item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    "product added to cart",
		"cart_item": mapCartItem(cartItem),
	})
}

// GET CART ITEMS

func GetCartItems(c *gin.Context) {

	userIDVal, exists := c.Get("userID")
	userID, ok := userIDVal.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var cartItems []models.CartItem
	if err := config.DB.
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&cartItems).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cart items"})
		return
	}

	resp := make([]CartItemResponse, 0, len(cartItems))
	for _, item := range cartItems {
		resp = append(resp, mapCartItem(item))
	}

	c.JSON(http.StatusOK, gin.H{"cart_items": resp})
}


//  UPDATE CART ITEM 


func UpdateCartItem(c *gin.Context) {

	userIDVal, exists := c.Get("userID")
	userID, ok := userIDVal.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cartItem models.CartItem
	if err := config.DB.
		Preload("Product").
		Where("id = ? AND user_id = ?", cartID, userID).
		First(&cartItem).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	if input.Quantity > int(cartItem.Product.StockQuantity) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock available"})
		return
	}

	cartItem.Quantity = input.Quantity

	if err := config.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "cart item updated",
		"cart_item": mapCartItem(cartItem),
	})
}

//DELETE CART ITEM 

func DeleteCartItem(c *gin.Context) {

	userIDVal, exists := c.Get("userID")
	userID, ok := userIDVal.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}

	result := config.DB.
		Where("id = ? AND user_id = ?", cartID, userID).
		Delete(&models.CartItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove cart item"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart item removed"})
}
