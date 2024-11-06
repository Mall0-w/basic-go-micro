package repository

import (
	. "authentication-service/models"
)

// UserRepository defines the interface for user data operations
type AuthRepository interface {
	FindUserByID(id uint) (*User, error)
	FindUserByEmail(email string) (*User, error)
	UserExistsByEmail(email string) (bool, error)
	FindTokenByUserID(id uint) (*RefreshToken, error)
	CreateNewRefreshToken(t *RefreshToken) error
	RevokeAllTokensByUserID(userId uint) error
}
