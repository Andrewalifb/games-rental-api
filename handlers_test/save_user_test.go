package handlerstest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Andrewalifb/games-rental-api/models"

	"github.com/gin-gonic/gin"
)

func TestSaveUser(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a mock request
	req, err := http.NewRequest("POST", "/api/v1/users/register", nil)
	if err != nil {
			t.Fatal(err)
	}

	// Create a router instance
	router := gin.Default()

	// Mock the handler function
	router.POST("/api/v1/users/register", func(c *gin.Context) {
			// Mock user data
			user := &models.User{
					ID:            1,
					FullName:      "",
					Email:         "",
					Password:      "",
					PhoneNumber:   "",
					Address:       "",
					DepositAmount: 0,
					Role:          "user",
					CreatedAt:     time.Time{},
					UpdatedAt:     time.Time{},
			}

			// Return success response
			c.JSON(http.StatusCreated, gin.H{"error": models.ErrorResponse{
					Code:    201,
					Message: fmt.Sprintf("Success created a new user with ID : %d", user.ID),
			},
					"data": user,
			})
	})

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
			t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Parse the actual response

var actual map[string]interface{}
json.Unmarshal(rr.Body.Bytes(), &actual)

// Define the expected response
expected := map[string]interface{}{
    "error": map[string]interface{}{
        "code":    201,
        "message": "Success created a new user with ID : 1",
    },
    "data": map[string]interface{}{
        "id":             1,
        "full_name":      "",
        "email":          "",
        "password":       "",
        "phone_number":   "",
        "address":        "",
        "deposit_amount": 0,
        "role":           "user",
        "created_at":     "0001-01-01T00:00:00Z",
        "updated_at":     "0001-01-01T00:00:00Z",
    },
}

// Compare the actual and expected responses
if !reflect.DeepEqual(actual, expected) {
    t.Errorf("Handler returned unexpected body: got %v want %v", actual, expected)
}


}
