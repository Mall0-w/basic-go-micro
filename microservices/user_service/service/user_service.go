package userservice

import (
	"errors"
	"fmt"
	"net/http"
	"user-service/dtos"
	e "user-service/errors"
	"user-service/models"
	"user-service/repository" // Assuming you'll have a repository layer

	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) *UserService {
	if userRepo == nil {
		userRepo = repository.NewMysqlUserRepository(nil)
	}
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id int64) (*dtos.User, *e.Error) {
	if id <= 0 {
		return nil, e.NewError(http.StatusBadRequest, "invalid Id", e.ErrInvalidUserData)
	}

	user, err := s.userRepo.FindByID(id)

	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusNotFound, "User Doesn't exist", e.ErrNotFound)
		}
		return nil, e.NewError(http.StatusInternalServerError, "Failed to get user", err)
	}

	return user.ToUserDTO(), nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(u *dtos.UserCreate) (*dtos.User, *e.Error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to create password", err)
	}

	exists, err := s.userRepo.ExistsByEmail(u.Email)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "error when looking up email", err)
	}

	if exists {
		return nil, e.NewError(http.StatusBadRequest, "Email already exists", e.ErrUserExists)
	}

	user := models.NewUser(u.Name, u.Email, string(hashPassword))
	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to create password", err)
	}

	return createdUser.ToUserDTO(), nil
}

// // UpdateUser updates an existing user
func (s *UserService) UpdateUser(u *dtos.User) (*dtos.User, *e.Error) {
	if u.Id <= 0 {
		return nil, e.NewError(http.StatusBadRequest, "Invalid Id", e.ErrInvalidUserData)
	}

	// Check if user exists
	existingUser, err := s.userRepo.FindByID(u.Id)
	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusNotFound, "User doesn't exist", e.ErrRecordNotFound)
		}
		return nil, e.NewError(http.StatusInternalServerError, "Failed to get user", err)
	}

	// Check if email is being changed and if it's already taken
	if u.Email != existingUser.Email {
		exists, err := s.userRepo.ExistsByEmail(u.Email)
		if err != nil {
			return nil, e.NewError(http.StatusInternalServerError, "failed to check user existence", err)
		}
		if exists {
			return nil, e.NewError(http.StatusBadRequest, "User with email already exists", e.ErrUserExists)
		}
	}

	existingUser.Name = u.Name
	existingUser.Email = u.Email

	updatedUser, err := s.userRepo.Update(existingUser)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to update user", err)
	}

	return updatedUser.ToUserDTO(), nil
}

// // DeleteUser removes a user by their ID
func (s *UserService) DeleteUser(id int64) *e.Error {
	if id <= 0 {
		return e.NewError(http.StatusBadRequest, "invalid Id", e.ErrInvalidUserData)
	}

	// Check if user exists
	exists, err := s.userRepo.ExistsByID(id)
	if err != nil {
		return e.NewError(http.StatusInternalServerError, "Failed to check for user", err)
	}
	if !exists {
		return e.NewError(http.StatusNotFound, "User doesn't exist", e.ErrRecordNotFound)
	}

	if err := s.userRepo.Delete(id); err != nil {
		return e.NewError(http.StatusInternalServerError, "Failed to delete user", err)
	}

	return nil
}

// validateUserData validates user input data
func (s *UserService) validateUserData(name, email string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", e.ErrInvalidUserData)
	}
	if email == "" {
		return fmt.Errorf("%w: email is required", e.ErrInvalidUserData)
	}
	// Add more validation as needed (e.g., email format, name length, etc.)
	return nil
}
