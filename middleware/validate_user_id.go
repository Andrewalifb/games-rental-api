package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

//router.PUT("/books/:id", CheckBookExists(), UpdateBook)

func GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Get user id from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message:  "User ID not found",
			}})
			c.Abort()
			return
		}

		// check if user id exists
		record := config.DB.Where("id = ?", userID).First(&user)
		if record.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: "User not found " + record.Error.Error(),
			}})
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", &user)
		c.Next()
	}
}