package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/controller"
	userservice "user-service/service"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Switch to test mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	userService := userservice.NewUserService(nil)
	userController := controller.NewUserController(userService)
	userController.DefineRoutes(r)
	return r
}

func TestHelloWorld(t *testing.T) {
	// Create a new router instance
	router := setupRouter()

	// Create a new HTTP recorder
	w := httptest.NewRecorder()

	// Create a new request
	req, err := http.NewRequest("GET", "/users/", nil)
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

	expectedResponse := "This is the users service"
	if response != expectedResponse {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponse, response)
	}
}
