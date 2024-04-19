package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// func CheckoutCart(c *gin.Context) {
// 	err := godotenv.Load("config/.env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %s", err)
// 	}

// 	daysLeft := os.Getenv("DEFAULT_DAYS_LEFT")

// 	rentalID := c.Param("rental_id")

// 	var items []models.Cart
// 	config.DB.Where("rental_id = ?", rentalID).Find(&items)

// 	var totalAmount float32
// 	for _, item := range items {
// 		totalAmount += item.Price * float32(item.Quantity)
// 	}

// 	var user models.User
// 	config.DB.First(&user, items[0].UserID)

// 	if user.DepositAmount < totalAmount {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient deposit amount"})
// 		return
// 	}

// 	config.DB.Model(&models.PaymentHistory{}).Where("rental_id = ?", rentalID).Update("status", "completed")

// 	for _, item := range items {
// 		rentalTransaction := models.RentalTransaction{
// 			RentalID:        item.RentalID,
// 			UserID:          item.UserID,
// 			GameID:          item.GameID,
// 			Quantity:        item.Quantity,
// 			Price:           item.Price,
// 			TotalRentalCost: item.Price * float32(item.Quantity),
// 			RentedAt:        time.Now(),
// 		}
// 		config.DB.Create(&rentalTransaction)

// 		config.DB.Model(&models.Game{}).Where("id = ?", item.GameID).UpdateColumn("availability", gorm.Expr("availability - ?", item.Quantity))
// 	}

// 	days, err := strconv.Atoi(daysLeft)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Cant create daysleft"})
// 		return
// 	}

// 	rent_maintenance := models.RentMaintenance{
// 		RentalID: rentalID,
// 		UserID: items[0].UserID,
// 		DaysLeft: days,
// 		Status: "",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}

// 	config.DB.Create(&rent_maintenance)

// 	config.DB.Delete(&models.Cart{}, "rental_id = ?", rentalID)

// 	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
// }

func ProcessCheckout(c *gin.Context) {
  // Get daysleft from .env
	daysLeft := os.Getenv("DEFAULT_DAYS_LEFT")
  // Get data from context
	rentalID := c.Param("rental_id")
	items := c.MustGet("items").([]models.Cart)

	var user models.User
	// Get user id data 
	if err := config.DB.First(&user, items[0].UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
			Code: 404,
			Message: err.Error(),
			Details: "User ID not found",
		}})
		return
	}
	

  // Count totalAmount of each items sub total for totl invoice
	var totalAmount float32
	for _, item := range items {
		totalAmount += item.Price * float32(item.Quantity)
	}
  
	// Check first if the user deposit balance is biggor or at least equal to the total invoice amount
	if user.DepositAmount < totalAmount {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: "Top Up your deposit balance",
			Details: "Insufficient deposit amount",
		}})
		return
	}

	// Update the payment history and retrieve the updated record
	var paymentHistory models.PaymentHistory
	if err := config.DB.Model(&paymentHistory).Where("rental_id = ?", rentalID).Update("status", "completed").First(&paymentHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: fmt.Sprintf("Error while update the status payment history of rental id %s", rentalID),
		}})
		return
	}

	// Iterate throug Cart items slice to prepare data for rental transaction table 
	for _, item := range items {
		rentalTransaction := models.RentalTransaction{
			RentalID:        item.RentalID,
			UserID:          item.UserID,
			GameID:          item.GameID,
			PaymentID:       paymentHistory.ID, // Add the PaymentID to the RentalTransaction
			Quantity:        item.Quantity,
			Price:           item.Price,
			TotalRentalCost: item.Price * float32(item.Quantity),
			RentedAt:        time.Now(),
		}

		// Insert on by one data into rental transaction database
    if err := config.DB.Create(&rentalTransaction).Error; err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: "Error insert data into rental transaction database",
			}})
	    return
     }
		// Reduce an item availability on games database to make sure the stock avaliable data real time
		
		if err := config.DB.Model(&models.Game{}).Where("id = ?", item.GameID).UpdateColumn("availability", gorm.Expr("availability - ?", item.Quantity)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: fmt.Sprintf("Error while update the availability of game ID %d", item.GameID),
			}})
			return
		}
	
	}

  // Reduce the user's deposit amount by the total amount of invoice
	user.DepositAmount -= totalAmount
	
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error while reduce the use Deposit Amount",
		}})
		return
	}

	// Send Success payment report to customer email
  EmailSuccessPayment(c, items)

	// convert daysleft from string into int
	days, err := strconv.Atoi(daysLeft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
		Code: 500,
		Message: err.Error(),
		Details: "Error when convert the days left to integer",
	}})
		return
	}

	// Insert data into rent maintenance database 
	rent_maintenance := models.RentMaintenance{
		RentalID: rentalID,
		UserID: items[0].UserID,
		DaysLeft: days,
		Status: "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&rent_maintenance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error insert data into rental maintennce database",
		}})
		return
	}

	// after all transaction checkout finish remove all rental_id related data from cart
	
  if err := config.DB.Where("rental_id = ?", rentalID).Delete(&models.Cart{}).Error; err != nil {
  	c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
  		Code: 500,
  		Message: err.Error(),
  		Details: fmt.Sprintf("Error delete rental_id %s from cart database", rentalID),
  	}})
		return 
  }


	c.JSON(http.StatusCreated, gin.H{"success": models.SuccessResponse{
		Code: 201,
		Message: fmt.Sprintf("Success check out and payment rental id %s", rentalID),
		Details: "Success payment detail has been sent to your email",
	}})
}
