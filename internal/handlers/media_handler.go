package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"go-rest/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadItemMedia(c *gin.Context) {
	itemID := c.Param("id")

	// Check if item exists
	var item models.Item
	if err := database.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Determine type based on content type
	mediaType := "image"
	if header.Header.Get("Content-Type") == "video/mp4" { // Simple check
		mediaType = "video"
	}

	url, publicID, err := services.UploadToCloudinary(file, "inventory/items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to Cloudinary"})
		return
	}

	media := models.Media{
		ItemID:   item.ID,
		Type:     mediaType,
		URL:      url,
		PublicID: publicID,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}
