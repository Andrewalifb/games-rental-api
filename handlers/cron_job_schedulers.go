package handlers

import (
	"log"
	"time"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/go-co-op/gocron"
)

func RegisterCronJob() {
	// Load the WIB timezone
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Error loading location: %s", err)
	}

	// Create a new scheduler using the WIB timezone
	s := gocron.NewScheduler(location)

	_, err = s.Every(1).Day().At("00:01").Do(func() {
		var rents []models.RentMaintenance
		config.DB.Where("days_left != ? AND status = ?", 0, "not returned").Find(&rents)

		for _, rent := range rents {
			rent.DaysLeft -= 1
			rent.UpdatedAt = time.Now()
			config.DB.Save(&rent)
		}
	})

	if err != nil {
		log.Fatalf("Error scheduling job: %s", err)
	}

	s.StartAsync()
}