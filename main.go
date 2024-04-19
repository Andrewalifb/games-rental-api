package main

import (
	"log"
	"os"

	"github.com/Andrewalifb/games-rental-api/config"
	"github.com/Andrewalifb/games-rental-api/handlers"
	"github.com/Andrewalifb/games-rental-api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting application ...")

  // Load .env file
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	
	// Create Database Connection
	config.DatabaseConnection()
	// Database Migration
	config.Migrate()

	// Register cron job
	handlers.RegisterCronJob()
	router := gin.Default()

   // api v1 routes
   v1 := router.Group("/api/v1") 
   {
   	users := v1.Group("/users")
   	{
   		users.POST("/register", middleware.Register() ,handlers.SaveUser) 
   		users.POST("/login", middleware.Login(), handlers.GetUserByEmail, handlers.LoginUser) 
   	}
   
   	auth := v1.Group("/")
   	auth.Use(middleware.Auth())
   	{
   		users := auth.Group("/users")
   		{
   			rentals := users.Group("/rentals") 
   			{
   				rentals.GET("/latest_rent", middleware.FetchLatestRentalID(), handlers.GetUserLatestRent)
   				rentals.GET("/all", handlers.GetUserAllRent)
   			}
   				
   			deposits := users.Group("/deposits") 
   			{
   				deposits.PUT("/topup", middleware.TopUp(), middleware.GetUserByID(), handlers.TopUpBalance)
   			}
   				
   			carts := users.Group("/carts") 
   			{
   				carts.GET("/", handlers.GetCartByUserID) 
   				carts.POST("/", middleware.BindCartItems() ,middleware.ValidateGameIDs(), handlers.AddItemsToCart) 
   				carts.POST("/checkout/:rental_id", middleware.CheckItemQuantity(), handlers.ProcessCheckout)
   			}
   		}
			  admin := auth.Group("/admin")
				{
					generate := admin.Group("/generate")
					{
						generate.GET("/bussiness_revenue", middleware.RoleValidation("admin"), handlers.GenerateBussinessRevenue)
						generate.GET("/top_customers", middleware.RoleValidation("admin"), handlers.GenerateTopCustomerList)
						generate.GET("/top_games", middleware.RoleValidation("admin"), handlers.GenerateTopGamesList)
					}
				}
   	}
   }

  // router.Run(":8080")

	port := os.Getenv("PORT")
  if port == "" {
	  log.Fatal("$PORT must be set")
	}
	router.Run(":"+port)
}