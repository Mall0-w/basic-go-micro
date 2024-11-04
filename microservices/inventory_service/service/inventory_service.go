package service

import (
	. "inventory-service/dtos"
	"inventory-service/repository"
)

type InventoryService struct {
	inventoryRepository *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepository: repo,
	}
}

func (is *InventoryService) GetInventoryById(id int64) (*Inventory, error) {

	return NewInventory(id, "test inventory"), nil
}
