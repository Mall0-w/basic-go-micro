package userservice

import (
	"errors"
	"fmt"
	. "user-service/dtos"
	"user-service/repository" // Assuming you'll have a repository layer
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidUserData = errors.New("invalid user data")
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) *UserService {

	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id int64) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid ID", ErrInvalidUserData)
	}

	// user, err := s.userRepo.FindByID(id)
	user := NewUser(id, "test user", "test user")
	var err error = nil

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// CreateUser creates a new user
// func (s *UserService) CreateUser(name, email string) (*User, error) {
// 	if err := s.validateUserData(name, email); err != nil {
// 		return nil, err
// 	}

// 	// Check if user with email already exists
// 	exists, err := s.userRepo.ExistsByEmail(email)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to check user existence: %w", err)
// 	}
// 	if exists {
// 		return nil, fmt.Errorf("%w: email=%s", ErrUserExists, email)
// 	}

// 	user := &User{
// 		Name:  name,
// 		Email: email,
// 	}

// 	createdUser, err := s.userRepo.Create(user)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create user: %w", err)
// 	}

// 	return createdUser, nil
// }

// // UpdateUser updates an existing user
// func (s *UserService) UpdateUser(id int64, name, email string) (*User, error) {
// 	if id <= 0 {
// 		return nil, fmt.Errorf("%w: invalid ID", ErrInvalidUserData)
// 	}

// 	if err := s.validateUserData(name, email); err != nil {
// 		return nil, err
// 	}

// 	// Check if user exists
// 	existingUser, err := s.userRepo.FindByID(id)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrNotFound) {
// 			return nil, fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
// 		}
// 		return nil, fmt.Errorf("failed to get user: %w", err)
// 	}

// 	// Check if email is being changed and if it's already taken
// 	if email != existingUser.Email {
// 		exists, err := s.userRepo.ExistsByEmail(email)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to check user existence: %w", err)
// 		}
// 		if exists {
// 			return nil, fmt.Errorf("%w: email=%s", ErrUserExists, email)
// 		}
// 	}

// 	existingUser.Name = name
// 	existingUser.Email = email

// 	updatedUser, err := s.userRepo.Update(existingUser)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update user: %w", err)
// 	}

// 	return updatedUser, nil
// }

// // DeleteUser removes a user by their ID
// func (s *UserService) DeleteUser(id int64) error {
// 	if id <= 0 {
// 		return fmt.Errorf("%w: invalid ID", ErrInvalidUserData)
// 	}

// 	// Check if user exists
// 	exists, err := s.userRepo.ExistsByID(id)
// 	if err != nil {
// 		return fmt.Errorf("failed to check user existence: %w", err)
// 	}
// 	if !exists {
// 		return fmt.Errorf("%w: id=%d", ErrUserNotFound, id)
// 	}

// 	if err := s.userRepo.Delete(id); err != nil {
// 		return fmt.Errorf("failed to delete user: %w", err)
// 	}

// 	return nil
// }

// validateUserData validates user input data
func (s *UserService) validateUserData(name, email string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidUserData)
	}
	if email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidUserData)
	}
	// Add more validation as needed (e.g., email format, name length, etc.)
	return nil
}
