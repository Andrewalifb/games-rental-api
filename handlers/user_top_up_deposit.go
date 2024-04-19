package handlers

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// {
// 	"amount": 100.00
// }

// func TopUpBalance(c *gin.Context) {
// 	var user models.User
// 	// id := c.Param("id")

// 	// Get user id from context
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "User ID not found",
// 		})
// 		return
// 	}

// 	var topUpAmount struct {
// 			Amount float32 `json:"amount"`
// 	}

// 	err := c.ShouldBindJSON(&topUpAmount)

// 	if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 					"error": err,
// 			})
// 			return
// 	}

// 	res := config.DB.Model(&user).Where("id = ?", userID).UpdateColumn("deposit_amount", gorm.Expr("deposit_amount + ?", topUpAmount.Amount))

// 	if res.RowsAffected == 0 {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 					"error": "balance not updated",
// 			})
// 			return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 			"message": "balance updated successfully",
// 	})

// }

func TopUpBalance(c *gin.Context) {
	// get data from context
	user, _ := c.Get("user")
	topUpAmount, _ := c.Get("topUpAmount")

	res := config.DB.Model(user.(*models.User)).UpdateColumn("deposit_amount", gorm.Expr("deposit_amount + ?", topUpAmount.(*struct{ Amount float32 "json:\"amount\"" }).Amount))

	if res.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: "Error user balance not updated",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": models.SuccessResponse{
		Code: 200,
		Message: "user deposit balance updated successfully",
	}})
}