package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func TopUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var topUpAmount struct {
			Amount float32 `json:"amount"`
		}

		err := c.ShouldBindJSON(&topUpAmount)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
			}})
			c.Abort()
			return
		}

		// Store topUpAmount in context
		c.Set("topUpAmount", &topUpAmount)
		c.Next()
	}
}