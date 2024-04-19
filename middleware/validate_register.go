package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Binding request body
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: "Error on binding user body " + err.Error(),
			}})
			c.Abort()
			return
		}
	  
		// Email validation
		if !IsValidEmail(user.Email) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: "Error Email is not valid",
			}})
			c.Abort()
			return
		}
	  
		// Hashing user password
		if err := user.HashPassword(user.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: "Error on hashing user password " + err.Error(),
			}})
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", &user)
		c.Next()
	}
}