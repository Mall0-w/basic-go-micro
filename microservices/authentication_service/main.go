package main

import (
	"authentication-service/controller"
	s "authentication-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Create service with repository
	userService := s.NewAuthService(nil)

	// Create controller with service
	userController := controller.NewAuthController(userService)
	userController.DefineRoutes(r)

	r.Run("0.0.0.0:8080")
}
