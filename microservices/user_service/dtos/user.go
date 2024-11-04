package dtos

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"user_id"`
	Password string `json:"password_hash"`
}

func NewUser(id int64, userID, passwordHash string) *User {
	return &User{
		ID:       id,
		Email:    userID,
		Password: passwordHash,
	}
}
