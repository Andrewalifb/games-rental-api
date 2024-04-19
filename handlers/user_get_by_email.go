package handlers

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func GetUserByEmail(c *gin.Context) {

	// Get token email request
	request, _ := c.Get("request")
	var user models.User

	// check if email exists and password is correct
	record := config.DB.Where("email = ?", request.(*models.TokenRequest).Email).First(&user)
	if record.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
			Code: 404,
			Message: "Error email not found",
		}})
		c.Abort()
		return
	}

	credentialError := user.CheckPassword(request.(*models.TokenRequest).Password)
	if credentialError != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrorResponse{
			Code: 401,
			Message: "Invalid password",
		}})
		c.Abort()
		return
	}

	// Store user in context
	c.Set("user", &user)
	c.Next()
}