package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetCartByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// Mock the JWT middleware
	r.Use(func(c *gin.Context) {
		c.Set("userID", "testUserID")
		c.Next()
	})

	r.GET("/", GetCartByUserID)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Please replace this with your expected response
	expected := `{"Success": {"Code": 200, "Message": "Success Get Cart Data by User ID"}, "Details": [{"id": 1, "rental_id": "1", "full_name": "John Doe", "game_name": "Game 1", "quantity": 1, "price": 10.0, "sub_total": 10.0, "added_at": "2024-04-18T17:34:18Z"}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
