package tests

import (
	"encoding/json"
	"inventory-service/controller"
	"inventory-service/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	userService := service.NewInventoryService(nil)
	userController := controller.NewInventoryController(userService)
	userController.DefineRoutes(r)
	return r
}

func TestHelloWorld(t *testing.T) {
	// Create a new router instance
	router := setupRouter()

	// Create a new HTTP recorder
	w := httptest.NewRecorder()

	// Create a new request
	req, err := http.NewRequest("GET", "/inventory/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Serve the request using our router
	router.ServeHTTP(w, req)

	// Check if the status code is correct
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	var response string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Couldn't parse response body: %v\n", err)
	}

	expectedResponse := "This is the Inventory Service"
	if response != expectedResponse {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponse, response)
	}
}
