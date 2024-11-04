package main

import (
	"inventory-service/controller"
	"inventory-service/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Create service with repository
	userService := service.NewInventoryService(nil)

	// Create controller with service
	userController := controller.NewInventoryController(userService)
	userController.DefineRoutes(r)

	r.Run("0.0.0.0:8080")
}
