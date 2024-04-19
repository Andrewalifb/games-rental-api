package handlers

import (
	"net/http"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

func GenerateTopGamesList(c *gin.Context) {

	var result []struct {
		GameName        string
		PlatformName    string
		TotalRentalCost float32
		LastRentedAt    time.Time
	}

	err := config.DB.Table("rental_transactions").
		Select("games.name as game_name, game_platforms.name as platform_name, SUM(rental_transactions.total_rental_cost) as total_rental_cost, MAX(rental_transactions.rented_at) as last_rented_at").
		Joins("left join games on games.id = rental_transactions.game_id").
		Joins("left join game_platforms on game_platforms.id = games.platform_id").
		Group("games.id, game_platforms.id").
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
		Message: "Success Get Games With Total Cost And Last Rent",
	},
		"data": result,
	})
}