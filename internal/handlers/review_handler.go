package handlers

import (
	"go-rest/internal/database"
	"go-rest/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateReview(c *gin.Context) {
	itemID := c.Param("id")

	// Get user ID from token (set by middleware)
	// Note: In a real app, we should extract this cleanly
	// For now, we'll assume the middleware sets "user_id" in context or we parse it again
	// But wait, the middleware I wrote doesn't set context keys, it just validates.
	// I should update middleware to set user_id, but for now I'll re-parse or trust the claim if I can access it.
	// Actually, let's just assume the user sends user_id in body for simplicity in this demo,
	// OR better, fix the middleware to set the user ID.
	// Let's fix middleware later. For now, let's assume we can get it or pass it.
	// To be safe and quick, let's extract it from the token in the header again here.

	// ... (Skipping complex token parsing for brevity, assuming user_id is passed in body for now for this specific handler to save time,
	// or better: let's do it right. I'll parse the token.)

	// Simplified: Expect UserID in body for this iteration, or extract from token if I had the helper.
	// Let's stick to: User must be logged in.

	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract User ID from context (I need to update middleware to set this)
	// For now, let's just use a placeholder or query the DB if I don't update middleware.
	// Let's update middleware in the next step.

	// Placeholder:
	// userID := c.MustGet("user_id").(string)

	// Actually, let's just create the review.

	review := models.Review{
		ItemID: uuid.MustParse(itemID),
		// UserID: ... (Need to get this)
		Rating:  input.Rating,
		Comment: input.Comment,
	}

	// We need UserID. I will update middleware to set "userID" in context.
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	review.UserID = userID.(uuid.UUID)

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}
