package dtos

type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func NewUser(id int64, name, email string) *User {
	return &User{
		Id:    id,
		Name:  name,
		Email: email,
	}
}
