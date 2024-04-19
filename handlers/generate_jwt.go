package handlers

import (
	"log"
	"os"
	"time"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/golang-jwt/jwt"

	"github.com/joho/godotenv"
)


func GenerateJWT(id uint, email, role string) (tokenString string, err error) {

	err = godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY")) 

  expirationTime := time.Now().Add(1 * time.Hour)
  claims:= &models.JWTClaim{
    UserID: id,
    Email: email,
    Role: role,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: expirationTime.Unix(),
    },
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  tokenString, err = token.SignedString(jwtKey)
  return
}

