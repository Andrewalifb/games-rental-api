package handlers

import (
	"net/http"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)



func GetCartByUserID(c *gin.Context) {

	type result struct {
		ID       uint      `json:"id"`
		RentalID string    `json:"rental_id"`
		FullName string    `json:"full_name"`
		GameName string    `json:"game_name"`
		Quantity int       `json:"quantity"`
		Price    float32   `json:"price"`
		SubTotal float32   `json:"sub_total"`
		AddedAt  time.Time `json:"added_at"`
	}
	
	// Get user id from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
			Code: 404,
			Message: "Error get User ID from context",
		}})
		return
	}

	var results []result

	// Database query using gorm to populate car data by user ID
	config.DB.Table("carts").Select("carts.id, carts.rental_id, users.full_name, games.name, carts.quantity, carts.price, carts.added_at").
		Joins("inner join users on users.id = carts.user_id").
		Joins("inner join games on games.id = carts.game_id").
		Where("carts.user_id = ?", userID).
		Order("carts.rental_id DESC").
		Scan(&results)

	// Check if result is empty
	if len(results) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
			Code: 404,
			Message: "Cart data not found!",
		}})
		return
	}

	// Calculate the SubTotal for each item
	for i, item := range results {
		results[i].SubTotal = float32(item.Quantity) * item.Price
	}

	c.JSON(http.StatusOK, gin.H{"Success": models.SuccessResponse{
		Code: 200,
		Message: "Success Get Cart Data by User ID",
	},
	"Details": results,
})
}