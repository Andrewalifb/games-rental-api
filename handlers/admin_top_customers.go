package handlers

import (
	"net/http"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func GenerateTopCustomerList(c *gin.Context) {

	var result []struct {
		UserID          uint
		FullName        string
		TotalRentalCost float32
		LastRentedAt    time.Time
	}

	err := config.DB.Table("rental_transactions").
		Select("users.id as user_id, users.full_name, SUM(rental_transactions.total_rental_cost) as total_rental_cost, MAX(rental_transactions.rented_at) as last_rented_at").
		Joins("left join users on users.id = rental_transactions.user_id").
		Group("users.id").
		Order("total_rental_cost desc").
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
			Code:    500,
			Message: "Failed to execute query",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": models.SuccessResponse{
		Code:    200,
		Message: "Success Get Users With Total Cost And Last Rent",
	},
		"data": result,
	})
}