package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type input struct {
	Name          string  `json:"name" binding:"required"`
	Category      string  `json:"category" binding:"required"`
	Description   string  `json:"description"`
	ImageUrl      string  `json:"image_url" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	StockQuantity int     `json:"stock_quantity" binding:"required"`
}

func CreateProduct(c *gin.Context) {
	var input input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existProduct models.Product

	if err:=config.DB.Where("name = ?",input.Name).First(&existProduct).Error;
	err == nil{
		c.JSON(http.StatusConflict,gin.H{"error":"this product already exist"})
		return
	}else if err != gorm.ErrRecordNotFound {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
	return
}
	

	product := models.Product{
		Name:          input.Name,
		Description:   input.Description,
		Category:      input.Category,
		ImageURL:      input.ImageUrl,
		Price:         input.Price,
		StockQuantity: uint(input.StockQuantity),
	}

	

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"product": product,
	})
}

func UpdateProduct(c *gin.Context) {
	IdParam := c.Param("id")
	id, err := strconv.ParseUint(IdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check your id"})
		return
	}

	var product models.Product

	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input input

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Description != "" {
		product.Description = input.Description
	}
	if input.Category != "" {
		product.Category = input.Category
	}
	if input.Price != 0 {
		product.Price = input.Price
	}
	if input.StockQuantity != 0 {
		product.StockQuantity = uint(input.StockQuantity)
	}
	if input.ImageUrl != "" {
		product.ImageURL = input.ImageUrl
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save into database"})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "product updated successfully",
		"product": product,
	})
}

func DeleteProduct(c *gin.Context) {
    // Get the product ID from the URL
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Delete the product
    res := config.DB.Delete(&models.Product{}, uint(id))
    if res.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
    }

    if res.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "Product deleted successfully"})
}

// get all products -----public----
func GetAllProducts(c *gin.Context) {
	var products []models.Product

	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot fetch the products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// get products By Id (public)
func GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "check your ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, uint(id)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})

}
