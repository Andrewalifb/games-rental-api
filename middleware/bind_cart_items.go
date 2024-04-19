package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func BindCartItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Cart

		// Binding JSON
		if err := c.ShouldBindJSON(&items); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
				Code: 400,
				Message: err.Error(),
			}})
			c.Abort()
			return
		}

		// Store items in context
		c.Set("items", items)
		c.Next()
	}
}

