package handlers

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"

	"github.com/gin-gonic/gin"
)

// {
// 	"email": "ardiandwinanda@mail.com",
// 	"password": "nanda_12345"
// }

// type TokenRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func LoginUser(context *gin.Context) {
// 	var request TokenRequest
// 	var user models.User

// 	// Binding request body
// 	if err := context.ShouldBindJSON(&request); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		context.Abort()
// 		return
// 	}

// 	// check if email exists and password is correct
// 	record := config.DB.Where("email = ?", request.Email).First(&user)
// 	if record.Error != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Email not found " + record.Error.Error()})
// 		context.Abort()
// 		return
// 	}

// 	credentialError := user.CheckPassword(request.Password)
// 	if credentialError != nil {
// 		context.JSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid Password"})
// 		context.Abort()
// 		return
// 	}

// 	// Generate JWT Token
// 	tokenString, err := GenerateJWT(user.ID, user.Email, user.Role)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error on generate JWT Token " + err.Error()})
// 		context.Abort()
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{
// 		"message": "Login Success",
// 		"token": tokenString})
// }


func LoginUser(c *gin.Context) {
	// Get User data from context
	user, _ := c.Get("user")

	// Generate JWT Token
	tokenString, err := GenerateJWT(user.(*models.User).ID, user.(*models.User).Email, user.(*models.User).Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error generating JWT Token",
		}})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": models.SuccessResponse{
		Code: 200,
		Message: "Login Success",
	},
	  "token": tokenString,
})
}


