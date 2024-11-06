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

// max length of a password, limited by bcrypt
var MAX_PASSWORD_LENGTH int = 70

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repository.UserRepository) *UserService {
	if userRepo == nil {
		//use default repo
		userRepo = repository.NewMysqlUserRepository(nil)
	}
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*dtos.User, *e.Error) {
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
func (s *UserService) CreateUser(u dtos.UserCreate) (*dtos.User, *e.Error) {

	//check that a user with the given email doesn't already exist
	exists, err := s.userRepo.ExistsByEmail(u.Email)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "error when looking up email", err)
	}

	if exists {
		return nil, e.NewError(http.StatusBadRequest, "Email already exists", e.ErrUserExists)
	}

	//bcrypt puts a cap on how long a password can be
	if len(u.Password) > MAX_PASSWORD_LENGTH {
		return nil, e.NewError(http.StatusBadRequest, "password too long", fmt.Errorf("password must be less than %d characters", MAX_PASSWORD_LENGTH))
	}

	//generate a salted hashed password with bcrypt
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to create password", err)
	}

	//create new user within db
	user := models.NewUser(u.Name, u.Email, string(hashPassword))
	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to create password", err)
	}

	return createdUser.ToUserDTO(), nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(u dtos.User) (*dtos.User, *e.Error) {
	if u.Id <= 0 {
		return nil, e.NewError(http.StatusBadRequest, "Invalid Id", e.ErrInvalidUserData)
	}

	//check if user exists
	existingUser, err := s.userRepo.FindByID(u.Id)
	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusNotFound, "User doesn't exist", e.ErrRecordNotFound)
		}
		return nil, e.NewError(http.StatusInternalServerError, "Failed to get user", err)
	}

	//check if email is being changed and if it's already taken
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

	//update info in db
	updatedUser, err := s.userRepo.Update(existingUser)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "failed to update user", err)
	}

	return updatedUser.ToUserDTO(), nil
}

// DeleteUser removes a user by their ID
func (s *UserService) DeleteUser(id uint) *e.Error {
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
