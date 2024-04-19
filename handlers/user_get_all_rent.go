package handlers

import (
	"net/http"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)


func GetUserAllRent(c *gin.Context) {

	var results []struct {
		RentalID        string
		FullName        string
		GameName        string
		PaymentID       uint
		Quantity        int
		Price           float32
		TotalRentalCost float32
		DaysLeft        int
		DueDate         time.Time
		RentStatus      string
	}

	// Get user id from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: "Error while try to get User ID from context",
		}})
		return
	}

	// Do database query using gorm 
	err := config.DB.Table("rental_transactions").
	Select("rental_transactions.rental_id, users.full_name, games.name as game_name, rental_transactions.payment_id, rental_transactions.quantity, rental_transactions.price, rental_transactions.total_rental_cost, rent_maintenances.days_left, (rental_transactions.rented_at + INTERVAL '1 day' * rent_maintenances.days_left) as due_date, rent_maintenances.status as rent_status").
	Joins("JOIN users ON rental_transactions.user_id = users.id").
	Joins("JOIN games ON rental_transactions.game_id = games.id").
	Joins("JOIN rent_maintenances ON rental_transactions.rental_id = rent_maintenances.rental_id").
	Where("rental_transactions.user_id = ?", userID).
	Order("rental_transactions.rented_at DESC").
	Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error getting rental transactions",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": models.ErrorResponse{
		Code: 200,
		Message: "Success get all user rent history",
	},
	"Details": results,})
}
