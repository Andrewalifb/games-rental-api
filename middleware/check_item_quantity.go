package middleware

import (
	"fmt"
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func CheckItemQuantity() gin.HandlerFunc {
	return func(c *gin.Context) {
		rentalID := c.Param("rental_id")

		// Check if user has add rental_id params
		if rentalID == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
				Code: 404,
				Message: "Invalid rental id or rental id not found",
			}})
			c.Abort()
			return
		}

		var items []models.Cart
		// Get cart daa by user id
		if err := config.DB.Where("rental_id = ?", rentalID).Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: "Database error occurred while fetching cart",
			}})
			c.Abort()
			return
		} else if len(items) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
				Code: 404,
				Message: fmt.Sprintf("No item found for the given rental id %s", rentalID),
			}})
			c.Abort()
			return
		}

		for _, item := range items {
			var game models.Game

			// Get game data by game id
			if err := config.DB.First(&game, item.GameID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
					Code: 500,
					Message: err.Error(),
					Details: "Database error occurred while fetching game",
				}})
				c.Abort()
				return
			} else if len(items) == 0 {
				c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
					Code: 404,
					Message: fmt.Sprintf("No items found for the given game ID %d", item.GameID),
				}})
				c.Abort()
				return
			}
      
			// Check if game item stock still avaliable for rent
			if game.Availability < item.Quantity {
				c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
					Code: 400,
					Message: fmt.Sprintf("Game ID %d is not available in the requested quantity", item.GameID),
				}})
				c.Abort()
				return
			}
		}

		c.Set("items", items)
		c.Next()
	}
}