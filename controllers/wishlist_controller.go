package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
func AddToWishlist(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var body struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var existing models.WishlistItem
	err := config.DB.
		Where("user_id = ? AND product_id = ?", userID, body.ProductID).
		First(&existing).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product already in wishlist"})
		return
	}

	wishlist := models.WishlistItem{
		UserID:    userID,
		ProductID: body.ProductID,
	}

	if err := config.DB.Create(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product into wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully added to wishlist"})
}

func GetWishlist(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var wishlist []models.WishlistItem
	if err := config.DB.
		Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&wishlist).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}
func RemoveFromWishlist(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case int:
		userID = uint(v)
	case uint:
		userID = v
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result := config.DB.
		Where("user_id = ? AND product_id = ?", userID, uint(productID)).
		Delete(&models.WishlistItem{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove product from wishlist"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found in wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from wishlist successfully"})
}
