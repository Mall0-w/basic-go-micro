package controller

import (
	"inventory-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	inventoryService *service.InventoryService
}

func NewInventoryController(inventoryService *service.InventoryService) *InventoryController {
	if inventoryService == nil {
		inventoryService = service.NewInventoryService(nil)
	}

	return &InventoryController{
		inventoryService: inventoryService,
	}
}

func (ic *InventoryController) DefineRoutes(r *gin.Engine) {
	inventoryGroup := r.Group("/inventory") // Using plural form as per REST conventions
	{
		inventoryGroup.GET("/", ic.testConnection)
		inventoryGroup.GET("/:id", ic.getInventoryById)
		// userGroup.POST("", uc.CreateUser)
		// userGroup.PUT("/:id", uc.UpdateUser)    // Added for completeness
		// userGroup.DELETE("/:id", uc.DeleteUser) // Added for completeness
	}
}

func (ic *InventoryController) testConnection(c *gin.Context) {
	c.JSON(http.StatusOK, "This is the Inventory Service")

	return
}

func (ic *InventoryController) getInventoryById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User ID",
		})
		return
	}

	inventory, err := ic.inventoryService.GetInventoryById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Inventory not found",
		})
		return
	}

	c.JSON(http.StatusOK, inventory)

}
