package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ValidateGameIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		items, _ := c.Get("items")

		for i, item := range items.([]models.Cart) {
			var game models.Game
			if err := config.DB.First(&game, item.GameID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{
						"error": fmt.Sprintf("Game ID %d not found", item.GameID),
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Database error",
					})
				}
				c.Abort()
				return
			}

			// Ensure the game's RentalCost is not empty
			if game.RentalCost == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
					Code: 401,
					Message: fmt.Sprintf("Game ID %d has an empty RentalCost", item.GameID),
				}})
				c.Abort()
				return
			}

			// Update the item's price with the game's RentalCost
			item.Price = game.RentalCost

			// Update the item in the slice
			items.([]models.Cart)[i] = item
		}

		// Store the updated items back in the context
		c.Set("items", items)

		c.Next()
	}
}



