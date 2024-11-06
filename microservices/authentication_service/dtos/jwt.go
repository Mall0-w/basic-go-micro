package dtos

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint   `json:"userID"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
