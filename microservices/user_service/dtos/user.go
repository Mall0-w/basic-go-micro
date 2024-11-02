package dtos

type User struct {
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}

func NewUser(id int64, userID, passwordHash string) *User {
	return &User{
		ID:           id,
		UserID:       userID,
		PasswordHash: passwordHash,
	}
}
