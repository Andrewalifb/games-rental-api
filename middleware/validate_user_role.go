package middleware

import (
	"net/http"

	"github.com/Andrewalifb/games-rental-api/models"
	"github.com/gin-gonic/gin"
)

// router.GET("/admin-endpoint", Auth(), RoleValidation("admin"), adminHandler)


func RoleValidation(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
				Code: 401,
				Message: "request does not contain an access token",
			}})
			c.Abort()
			return
		}
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrorResponse{
				Code: 401,
				Message: err.Error(),
			}})
			c.Abort()
			return
		}
		if claims.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": models.ErrorResponse{
				Code: 403,
				Message: "forbidden, not enough privileges",
			}})
			c.Abort()
			return
		}
		c.Next()
	}
}