package repository

import (
	. "user-service/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// function used to find a user given an id
	FindByID(id uint) (*User, error)
	// function used to see if a user exists by an id
	ExistsByID(id uint) (bool, error)
	//function used to see if a user exists by an email
	ExistsByEmail(email string) (bool, error)
	//function used to see create a user
	Create(user *User) (*User, error)
	//function used to update a user
	Update(user *User) (*User, error)
	//function used to delete a user
	Delete(id uint) error
}
