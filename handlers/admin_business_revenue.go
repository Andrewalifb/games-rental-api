package handlers

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func GenerateBussinessRevenue(c *gin.Context) {

	var result []struct {
			FullName            string
			GameName            string
			CategoryName        string
			PlatformName        string
			TotalRentalCost     float32
			PaymentAmount       float32
			PaymentStatus       string
			MaintenanceDaysLeft int
			MaintenanceStatus   string
	}

	totalRevenue := float32(0)

	err := config.DB.Table("rental_transactions").
			Select("users.full_name, games.name as game_name, game_categories.name as category_name, game_platforms.name as platform_name, rental_transactions.total_rental_cost, payment_histories.amount as payment_amount, payment_histories.status as payment_status, rent_maintenances.days_left as maintenance_days_left, rent_maintenances.status as maintenance_status").
			Joins("left join users on users.id = rental_transactions.user_id").
			Joins("left join games on games.id = rental_transactions.game_id").
			Joins("left join game_categories on game_categories.id = games.category_id").
			Joins("left join game_platforms on game_platforms.id = games.platform_id").
			Joins("left join payment_histories on payment_histories.rental_id = rental_transactions.rental_id").
			Joins("left join rent_maintenances on rent_maintenances.rental_id = rental_transactions.rental_id").
			Scan(&result).Error

	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": models.ErrorResponse{
					Code:    500,
					Message: err.Error(),
			}})
			return
	}

	for _, item := range result {
			totalRevenue += item.TotalRentalCost
	}

	c.JSON(http.StatusOK, gin.H{"success": models.SuccessResponse{
			Code: 200,
			Message: "Success Generate Bussiness Revenue",
	},
	"data":result,
	"totalRevenue": totalRevenue,
	})
}
