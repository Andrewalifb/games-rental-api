package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GetUserIdFromJWT(c *gin.Context) (string, error) {

  // Load JWT key from .env
	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	signedToken := c.GetHeader("Authorization")
	if signedToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
			Code: 401,
			Message: "request does not contain an access token",
		}})
		c.Abort()
		return "", errors.New("request does not contain an access token")
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*models.JWTClaim)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}

	return fmt.Sprintf("%d", claims.UserID), nil
}