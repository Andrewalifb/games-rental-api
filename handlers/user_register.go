package handlers

import (
	"fmt"
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"

	"github.com/Andrewalifb/games-rental-api/models"

	"github.com/gin-gonic/gin"
)

// Example json body request postman
// {
// 	"full_name": "Ardian Dwi Nanda",
// 	"email": "ardiandwinanda@mail.com",
// 	"password": "nanda_12345",
// 	"phone_number": "087712350987",
// 	"address": "Perumahan Bhayangkara Residence, Mambalan, Lombok Barat",
// 	"deposit_amount": 0,
// 	"role": ""
// }

// func RegisterUser(context *gin.Context) {
// 	var user models.User

// 	// Binding request body
// 	if err := context.ShouldBindJSON(&user); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Error on binding user body " + err.Error()})
// 		context.Abort()
// 		return
// 	}

// 	// Email validation
//   if !middleware.IsValidEmail(user.Email) {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Email is not valid"})
// 		context.Abort()
// 		return
// 	}

// 	// Hashing user password
// 	if err := user.HashPassword(user.Password); err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error on hashing user password " + err.Error()})
// 		context.Abort()
// 		return
// 	}

// 	// Insert user registration data into database
// 	record := config.DB.Create(&user)
// 	if record.Error != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error on creating register into database" + record.Error.Error()})
// 		context.Abort()
// 		return
// 	}

//    user.Password = ""
// 	// Return success response
// 	context.JSON(http.StatusCreated, gin.H{
// 		"message": fmt.Sprintf("Success created a new user with ID : %d", user.ID),
// 		"data": user,
// 	})

// }

func SaveUser(c *gin.Context) {
	// Get user data from context
	user, _ := c.Get("user")

	// Insert user registration data into database
	record := config.DB.Create(user.(*models.User))
	if record.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: "Error create user register data into database",
		}})
		return
	}

	// Clear the password before sending the response
	user.(*models.User).Password = ""

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"success": models.ErrorResponse{
		Code: 201,
		Message: fmt.Sprintf("Success created a new user with ID : %d", user.(*models.User).ID),
	},
   "data": user,
})
}