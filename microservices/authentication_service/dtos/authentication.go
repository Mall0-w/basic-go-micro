package dtos

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken  string
	RefreshToken string
}

func NewUserLoginResponse(accessToken, refreshToken string) *UserLoginResponse {
	return &UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

type RefreshResponse struct {
	AccessToken string
}

func NewRefreshResponse(accessToken string) *RefreshResponse {
	return &RefreshResponse{
		AccessToken: accessToken,
	}
}
