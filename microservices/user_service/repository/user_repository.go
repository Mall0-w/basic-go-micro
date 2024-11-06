package repository

import (
	. "user-service/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id uint) (*User, error)
	ExistsByID(id uint) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(id uint) error
}
