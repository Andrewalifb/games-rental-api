package middleware

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(signedToken string) (*models.JWTClaim, error) {

  jwtKey := []byte(os.Getenv("JWT_SECRET_KEY")) 

  // Validate token with claims
  token, err := jwt.ParseWithClaims(
      signedToken,
      &models.JWTClaim{},
      func(token *jwt.Token) (interface{}, error) {
          return []byte(jwtKey), nil
      },
  )

  if err != nil {
      return nil, err
  }

  // Get claims
  claims, ok := token.Claims.(*models.JWTClaim)
  if !ok {
      return nil, errors.New("couldn't parse claims")
  }

  // Check token expired time
  if claims.ExpiresAt < time.Now().Local().Unix() {
      return nil, errors.New("token expired")
  }
  return claims, nil
}

func Auth() gin.HandlerFunc {
  return func(c *gin.Context) {
      tokenString := c.GetHeader("Authorization")
      if tokenString == "" {
          c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrorResponse{
            Code: 401,
            Message: "request does not contain an access token",
          }})
          c.Abort()
          return
      }
      claims, err := ValidateToken(tokenString)
      if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": models.ErrorResponse{
          Code: 401,
          Message: "request does not contain an access token",
        }})
          c.Abort()
          return
      }
      // Store the user ID in the context
      c.Set("userID", claims.UserID)
      c.Next()
  }
}
