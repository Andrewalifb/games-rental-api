package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func CreateInvoice(c *gin.Context, rentalID int) (*models.Invoice, error) {
	// Fetch user details
	var user models.User

	items, _ := c.Get("items")
	userID, _ := c.Get("userID")
  
	// Load the third party api credential data from .env
	xenditApiKeys := os.Getenv("XENDIT_API_KEY")
	xenditApiUrl := os.Getenv("XENDIT_API_URL")

	// Get user data by user id
  if err := config.DB.First(&user, userID.(uint)).Error; err != nil {
  	c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
  		Code: 500,
  		Message: err.Error(),
  		Details: fmt.Sprintf("Error while get user id %v data from database", userID.(uint)),
  	}})
		return nil, err
  }

	// Set credentials to access third party api
	apiKey := xenditApiKeys
	apiUrl := xenditApiUrl

	// Prepare the items for the invoice
	var invoiceItems []interface{}
	var totalAmount float32

	// iterate through items of []model.Cart
	for _, item := range items.([]models.Cart) {
		var game models.Game

		// Get game name
		
		if err := config.DB.First(&game, item.GameID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": models.ErrorResponse{
				Code: 404,
				Message: err.Error(),
				Details: fmt.Sprintf("Error game_id %d is not found", item.GameID),
			}})
			return nil, err
		}

		// create each invoice item
		invoiceItem := map[string]interface{}{
			"name":     game.Name,
			"quantity": item.Quantity,
			"price":    item.Price,
		}

		invoiceItems = append(invoiceItems, invoiceItem)

		// Sum up the total rental cost
		totalAmount += item.Price * float32(item.Quantity)
	}

	// Define the available banks, retail outlets, and e-wallets according to Xendit Docummentaion
	availableBanks := []map[string]interface{}{
		{
			"bank_code": "MANDIRI",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
		},
		{
			"bank_code": "BRI",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "BNI",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "PERMATA",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "BCA",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "SAHABAT_SAMPOERNA",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "CIMB",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "BSI",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "BJB",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	{
			"bank_code": "BNC",
			"collection_type": "POOL",
			"transfer_amount": totalAmount,
			"bank_branch": "Virtual Account",
			"account_holder_name": "CARDS AUTOMATION TEST",
			"identity_amount": 0,
	},
	}
	availableRetailOutlets := []map[string]interface{}{
		{
			"retail_outlet_name": "ALFAMART",
		},
		{
			"retail_outlet_name": "INDOMARET",
		},
	}
	availableEwallets := []map[string]interface{}{
		{
			"ewallet_type": "SHOPEEPAY",
		},
		{
			"ewallet_type": "ASTRAPAY",
		},
		{
			"ewallet_type": "JENIUSPAY",
	  },
	  {
			"ewallet_type": "DANA",
	  },
	  {
			"ewallet_type": "LINKAJA",
	  },
	  {
			"ewallet_type": "OVO",
	  },
	  {
			"ewallet_type": "NEXCASH",
	  },
	}

	// Set the body request for Xendit invoice
	bodyRequest := map[string]interface{}{
		"external_id":      strconv.Itoa(rentalID),
		"amount":           totalAmount, // Use the total amount
		"description":      "Invoice for game rental",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":  user.FullName,
			"email": user.Email,
		},
		"customer_notification_preference": map[string]interface{}{
			"invoice_created": []string{
				"email",
			},
		},
		"currency":                 "IDR",
		"items":                    invoiceItems,
		"available_banks":          availableBanks,
		"available_retail_outlets": availableRetailOutlets,
		"available_ewallets":       availableEwallets,
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var resInvoice models.Invoice
	if err := json.NewDecoder(response.Body).Decode(&resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil
}
