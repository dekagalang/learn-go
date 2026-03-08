package handlers

import (
	"net/http"

	"crud-api/database"
	"crud-api/models"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []models.Product

	database.DB.Preload("User").Find(&products)

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	if err := database.DB.Preload("User").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	database.DB.Create(&product)

	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}

	c.ShouldBindJSON(&product)

	database.DB.Save(&product)

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	database.DB.Delete(&models.Product{}, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted",
	})
}
