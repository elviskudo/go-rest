package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ToggleFavorite(c *gin.Context) {
	itemID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	uid := userID.(uuid.UUID)
	iid := uuid.MustParse(itemID)

	var favorite models.Favorite
	if err := database.DB.Where("user_id = ? AND item_id = ?", uid, iid).First(&favorite).Error; err != nil {
		// Not favorited, so create it
		newFav := models.Favorite{
			UserID: uid,
			ItemID: iid,
		}
		database.DB.Create(&newFav)

		// Increment count
		database.DB.Model(&models.Item{}).Where("id = ?", iid).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))

		c.JSON(http.StatusCreated, gin.H{"message": "Favorited"})
	} else {
		// Favorited, so delete it
		database.DB.Delete(&favorite)

		// Decrement count
		database.DB.Model(&models.Item{}).Where("id = ?", iid).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))

		c.JSON(http.StatusOK, gin.H{"message": "Unfavorited"})
	}
}
