package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
// const CONFIG_SENDER_NAME = "Rental Games Store"
// const CONFIG_AUTH_EMAIL = "andrewalifb@gmail.com"
// const CONFIG_AUTH_PASSWORD = "yjmd dwws acde yazm"

func EmailSuccessPayment(c *gin.Context, items []models.Cart) {
	CONFIG_SENDER_NAME := os.Getenv("CONFIG_SENDER_NAME")
	CONFIG_AUTH_EMAIL := os.Getenv("CONFIG_AUTH_EMAIL")
	CONFIG_AUTH_PASSWORD := os.Getenv("CONFIG_AUTH_PASSWORD")
	// Fetch user details
	var user models.User
	if err := config.DB.First(&user, items[0].UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
			Code: 404,
			Message: err.Error(),
			Details: fmt.Sprintf("Error user id %d not found", items[0].UserID),
		}})
		return
	}

	// Prepare the items for the email
	var emailItems string
	var totalAmount float32

	// iterate through cart items for email items
	for _, item := range items {
		var game models.Game
		if err := config.DB.First(&game, item.GameID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
				Code: 404,
				Message: err.Error(),
				Details: fmt.Sprintf("Error game id %d not found", item.GameID),
			}})
			return
		}

		emailItem := fmt.Sprintf("Name: %s, Quantity: %d, Price: %.2f<br>", game.Name, item.Quantity, item.Price)
		emailItems += emailItem

		// Sum up the total rental cost
		totalAmount += item.Price * float32(item.Quantity)
	}

	subject := "Payment Successful"
	message := fmt.Sprintf("Dear %s,<br><br>Your payment for the following items has been successful:<br><br>%s<br><br>Total Amount: %.2f<br><br>Invoice ID: %s<br><br>Thank you for your purchase!<br><br>Best,<br>%s", user.FullName, emailItems, totalAmount, items[0].RentalID, CONFIG_SENDER_NAME)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", user.Email)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code: 500,
			Message: err.Error(),
			Details: "Error occurred while sending email",
		}})
		return
	}

}

