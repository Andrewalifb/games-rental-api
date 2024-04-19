package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"

	"github.com/gin-gonic/gin"
)



func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.TokenRequest

		// Binding request body
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
			}})
			c.Abort()
			return
		}

		// Store request in context
		c.Set("request", &request)
		c.Next()
	}
}


