package repository

import (
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type InventoryRepository interface {
}
