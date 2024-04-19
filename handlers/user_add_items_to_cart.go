package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

// [
//     {
//         "game_id": 8,
//         "quantity": 2
//     },
//     {
//         "game_id": 9,
//         "quantity": 1
//     }
// ]

// func AddToCart(c *gin.Context) {
// 	var items []models.Cart

// 	// Get user id from context
// 	userIDInterface, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "User ID not found",
// 		})
// 		return
// 	}

// 	// Perform a type assertion to convert interface{} to uint
// 	userID, ok := userIDInterface.(uint)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Error converting user ID to uint",
// 		})
// 		return
// 	}

// 	// Binding JSON
// 	if err := c.ShouldBindJSON(&items); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Start a new transaction
// 	tx := config.DB.Begin()

// 	// Get a new RentalID
// 	var rentalID int
// 	row := tx.Raw("SELECT nextval('rental_id_seq')").Row()
// 	row.Scan(&rentalID)

// 	// Calculate the total amount
// 	var totalAmount float32
// 	for _, item := range items {
// 		totalAmount += item.Price * float32(item.Quantity)
// 	}

// 	// Create a new payment history
// 	paymentHistory := models.PaymentHistory{
// 		UserID:          userID,
// 		RentalID:        strconv.Itoa(rentalID),
// 		Amount:          totalAmount,
// 		Status:          "pending",
// 		TransactionDate: time.Now(),
// 	}

// 	if err := tx.Create(&paymentHistory).Error; err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Insert each item into the cart
// 	for _, item := range items {
// 		item.RentalID = strconv.Itoa(rentalID)
// 		if err := tx.Create(&item).Error; err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 	}

// 	// Commit the transaction
// 	tx.Commit()

// 	c.JSON(http.StatusOK, gin.H{"message": "Items added to cart and payment history created successfully"})
// }

func AddItemsToCart(c *gin.Context) {
	items, _ := c.Get("items")
	userID, _ := c.Get("userID")

	// Start a new transaction
	tx := config.DB.Begin()

	// Generate a new RentalID
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
		UserID:          userID.(uint),
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
			Details: "Error while insert data into payment history database",
		}})
		return
	}

	// Insert each item into the cart
	for _, item := range items.([]models.Cart) {
		item.RentalID = strconv.Itoa(rentalID)
		item.UserID = userID.(uint) // Set the user_id for the item
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
				Code: 500,
				Message: err.Error(),
				Details: "Error while insert data into cart",
			}})
			return
		}
	}

	// Call CreateInvoice function here
	invoice, err := CreateInvoice(c, rentalID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error while create invoice",
		}})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Success response
	c.JSON(http.StatusCreated, gin.H{"success": models.SuccessResponse{
		Code: 201,
		Message: "Success ad items to cart",
	}, "Detail": invoice})
}

