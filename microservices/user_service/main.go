package main

import (
	"user-service/controller"
	userservice "user-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Create service with repository
	userService := userservice.NewUserService(nil)

	// Create controller with service
	userController := controller.NewUserController(userService)
	userController.DefineRoutes(r)

	r.Run("0.0.0.0:8080")
}
