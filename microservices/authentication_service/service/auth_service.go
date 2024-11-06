package AuthService

import (
	c "authentication-service/config"
	"authentication-service/dtos"
	e "authentication-service/errors"
	"authentication-service/models"
	"authentication-service/repository" // Assuming you'll have a repository layer
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles business logic for user operations
type AuthService struct {
	AuthRepo repository.AuthRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(AuthRepo repository.AuthRepository) *AuthService {
	if AuthRepo == nil {
		AuthRepo = repository.NewMysqlAuthRepository(nil)
	}
	return &AuthService{
		AuthRepo: AuthRepo,
	}
}

func (s *AuthService) UserLogin(u *dtos.UserLogin) (*dtos.UserLoginResponse, *e.Error) {
	existing, err := s.AuthRepo.FindUserByEmail(u.Email)
	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusNotFound, "No user exists with matching email", err)
		}
		return nil, e.NewError(http.StatusInternalServerError, "An error occurred when fetching the user", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(u.Password)); err != nil {
		return nil, e.NewError(http.StatusUnauthorized, "Invalid password", err)
	}

	existingDTO := existing.ToUserDTO()

	//generate access token
	accessToken, err := s.generateJWT(existingDTO, 15*time.Minute)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "Failed to generate access token", err)
	}

	//generate refresh token
	rawRefreshToken, err := s.generateJWT(existingDTO, 7*24*time.Hour)
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "Failed to generate refresh token", err)
	}

	//hash refresh for verification
	hashedToken := s.hashToken(rawRefreshToken)

	refreshToken := &models.RefreshToken{
		UserID:    existing.ID,
		TokenHash: string(hashedToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	//store refresh token in database and revoke all old refresh tokens
	if err := s.AuthRepo.CreateNewRefreshToken(refreshToken); err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "Failed to store refresh token", err)
	}

	return dtos.NewUserLoginResponse(accessToken, rawRefreshToken), nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*dtos.RefreshResponse, *e.Error) {
	//parse and validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, e.NewError(http.StatusUnauthorized, "Unexpected signing method", nil)
		}
		return []byte(c.LoadConfig().JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, e.NewError(http.StatusUnauthorized, "Invalid refresh token", err)
	}

	// Extract userID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["userID"] == nil {
		return nil, e.NewError(http.StatusUnauthorized, "Invalid token claims", nil)
	}
	//cast user id
	userID := uint(claims["userID"].(float64))

	//ensure user actually exists
	user, err := s.AuthRepo.FindUserByID(userID)

	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusNotFound, "User Doesn't exist", e.ErrNotFound)
		}
		return nil, e.NewError(http.StatusInternalServerError, "Failed to get user", err)
	}

	//get token stored in database for comparison
	storedToken, err := s.AuthRepo.FindTokenByUserID(userID)

	if err != nil {
		if errors.Is(err, e.ErrRecordNotFound) {
			return nil, e.NewError(http.StatusUnauthorized, "no valid refresh token found in db", e.ErrNotFound)
		}
		return nil, e.NewError(http.StatusInternalServerError, "Failed to get token", err)
	}

	//ensure stored token lines up with given token
	if s.hashToken(refreshToken) != storedToken.TokenHash {
		return nil, e.NewError(http.StatusUnauthorized, "refresh token does not matched stored token", fmt.Errorf("refresh token does not matched stored token"))
	}

	//generate new access token
	newAccessToken, err := s.generateJWT(user.ToUserDTO(), 15*time.Minute) // short-lived access token
	if err != nil {
		return nil, e.NewError(http.StatusInternalServerError, "Failed to generate new access token", err)
	}

	return dtos.NewRefreshResponse(newAccessToken), nil
}

func (s *AuthService) Logout(userId uint) *e.Error {
	if err := s.AuthRepo.RevokeAllTokensByUserID(userId); err != nil {
		return e.NewError(http.StatusInternalServerError, "Failed to log out user", err)
	}
	return nil
}

// Helper function to generate JWT
func (s *AuthService) generateJWT(u *dtos.User, duration time.Duration) (string, error) {

	claims := dtos.CustomClaims{
		UserID: u.Id,
		Email:  u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.LoadConfig().JwtSecret)) // Sign with the secret
}

func (s *AuthService) ParseJWT(tokenString string) (*dtos.CustomClaims, error) {
	claims := &dtos.CustomClaims{}

	// Parse the token with the claims and validate
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.LoadConfig().JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
