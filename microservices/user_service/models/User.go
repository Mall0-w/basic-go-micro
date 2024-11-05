package models

import "user-service/dtos"

type User struct {
	ID       int64 `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (u *User) ToUserDTO() *dtos.User {
	return &dtos.User{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
