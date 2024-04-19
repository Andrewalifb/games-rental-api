package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func FetchLatestRentalID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID")

		// Get the latest RentalID for the user
		var latestRentalID string
		err := config.DB.Table("rental_transactions").
			Select("rental_id").
			Where("user_id = ?", userID).
			Order("rented_at DESC").
			Limit(1).
			Scan(&latestRentalID).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: "Error getting latest RentalID",
			}})
			c.Abort()
			return
		}

		c.Set("latestRentalID", latestRentalID)
		c.Next()
	}
}