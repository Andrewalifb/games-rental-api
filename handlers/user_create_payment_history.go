package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func CreatePaymentHistory(c *gin.Context)  {
	  // Get items data from context
		items, _ := c.Get("items")

		// Get user id from context
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
				Code: 404,
				Message: "Error while try toget userID from context",
			}})
			c.Abort()
			return
		}

		// Perform a type assertion to convert interface{} to uint
		userID, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: "Error while perform type assertion of userID",
			}})
			c.Abort()
			return
		}

		// Start a new transaction
		tx := config.DB.Begin()

		// Get a new RentalID
		var rentalID int
		row := tx.Raw("SELECT nextval('rental_id_seq')").Row()
		row.Scan(&rentalID)

		// Calculate the total amount
		var totalAmount float32
		for _, item := range items.([]models.Cart) {
			totalAmount += item.Price * float32(item.Quantity)
		}

		// Create a new payment history
		paymentHistory := models.PaymentHistory{
			UserID:          userID,
			RentalID:        strconv.Itoa(rentalID),
			Amount:          totalAmount,
			Status:          "pending",
			TransactionDate: time.Now(),
		}

		if err := tx.Create(&paymentHistory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: "Error while insert data into payment history",
			}})
			c.Abort()
			return
		}

		// Store rentalID in context
		c.Set("rentalID", rentalID)

		// Commit the transaction
		tx.Commit()
}