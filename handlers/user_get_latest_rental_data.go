package handlers

import (
	"net/http"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

// func GetUserLatestRent(c *gin.Context) {
// 	// Get user id from context
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "User ID not found",
// 		})
// 		return
// 	}

// 	var results []struct {
// 		RentalID        string
// 		FullName        string
// 		GameName        string
// 		PaymentID       uint
// 		Quantity        int
// 		Price           float32
// 		TotalRentalCost float32
// 		DaysLeft        int
// 		DueDate         time.Time
// 		RentStatus      string
// 	}

// 	// Get the latest RentalID for the user
// 	var latestRentalID string
// 	err := config.DB.Table("rental_transactions").
// 		Select("rental_id").
// 		Where("user_id = ?", userID).
// 		Order("rented_at DESC").
// 		Limit(1).
// 		Scan(&latestRentalID).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting latest RentalID: " + err.Error()})
// 		return
// 	}

// 	err = config.DB.Table("rental_transactions").
// 		Select("rental_transactions.rental_id, users.full_name, games.name as game_name, rental_transactions.payment_id, rental_transactions.quantity, rental_transactions.price, rental_transactions.total_rental_cost, rent_maintenances.days_left, DATE_ADD(rental_transactions.rented_at, INTERVAL rent_maintenances.days_left DAY) as due_date, rent_maintenances.status as rent_status").
// 		Joins("JOIN users ON rental_transactions.user_id = users.id").
// 		Joins("JOIN games ON rental_transactions.game_id = games.id").
// 		Joins("JOIN rent_maintenances ON rental_transactions.rental_id = rent_maintenances.rental_id").
// 		Where("rental_transactions.user_id = ? AND rental_transactions.rental_id = ?", userID, latestRentalID).
// 		Scan(&results).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting rental transactions: " + err.Error()})
// 		return
// 	}

// 	var totalAmount float32
// 	err = config.DB.Table("payment_histories").
// 		Select("amount").
// 		Where("rental_id = ?", latestRentalID).
// 		Scan(&totalAmount).Error
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting total amount: " + err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "This is your latest rent data",
// 		"data": results,
// 		"total_amount": totalAmount,
// 	})
// }

func GetUserLatestRent(c *gin.Context) {
	// Get data from context
	userID := c.MustGet("userID")
	latestRentalID := c.MustGet("latestRentalID").(string)

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

	err := config.DB.Table("rental_transactions").
	Select("rental_transactions.rental_id, users.full_name, games.name as game_name, rental_transactions.payment_id, rental_transactions.quantity, rental_transactions.price, rental_transactions.total_rental_cost, rent_maintenances.days_left, (rental_transactions.rented_at + INTERVAL '1 day' * rent_maintenances.days_left) as due_date, rent_maintenances.status as rent_status").
	Joins("JOIN users ON rental_transactions.user_id = users.id").
	Joins("JOIN games ON rental_transactions.game_id = games.id").
	Joins("JOIN rent_maintenances ON rental_transactions.rental_id = rent_maintenances.rental_id").
	Where("rental_transactions.user_id = ? AND rental_transactions.rental_id = ?", userID, latestRentalID).
	Scan(&results).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error to retrieve rental transactions data",
		}})
		return
	}

	var totalAmount float32
	err = config.DB.Table("payment_histories").
		Select("amount").
		Where("rental_id = ?", latestRentalID).
		Scan(&totalAmount).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error getting total amount from payment history database",
		}})
		return
	}


	c.JSON(http.StatusOK, gin.H{"Success": models.SuccessResponse{
		Code: 200,
		Message: "Success Retrieve User Latest Rent Data",
	},
	  "Data": results,
	  "Total_Amount": totalAmount,
  })  
}