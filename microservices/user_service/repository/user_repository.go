package repository

import (
	. "user-service/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	FindByID(id int64) (*User, error)
	ExistsByID(id int64) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	Delete(id int64) error
}
