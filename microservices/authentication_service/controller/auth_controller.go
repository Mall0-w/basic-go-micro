package controller

import (
	// "fmt"
	conf "authentication-service/config"
	"authentication-service/dtos"
	. "authentication-service/service"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthController handles user-related HTTP requests
type AuthController struct {
	AuthService *AuthService // Assuming you have a AuthService
}

// NewAuthController creates a new AuthController instance
func NewAuthController(AuthService *AuthService) *AuthController {
	if AuthService == nil {
		AuthService = NewAuthService(nil)
	}

	return &AuthController{
		AuthService: AuthService,
	}
}

func (ac *AuthController) DefineRoutes(r *gin.Engine) {
	userGroup := r.Group("/auth") // Using plural form as per REST conventions
	{
		userGroup.GET("/", ac.TestConnection)
		// userGroup.GET("/:id", uc.GetUserByID)
		userGroup.POST("/login", ac.LoginUser)
		userGroup.GET("/claims", ac.ShowClaims)
		userGroup.POST("/logout", ac.LogoutUser)
		userGroup.GET("/refresh", ac.RefreshToken)
	}
}

func (ac *AuthController) TestConnection(c *gin.Context) {
	c.JSON(http.StatusOK, "This is the auth service")
}

func (ac *AuthController) LoginUser(c *gin.Context) {
	var request dtos.UserLogin
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	response, e := ac.AuthService.UserLogin(&request)
	if e != nil {
		c.JSON(e.Code, e.ToJson())
		return
	}

	c.SetCookie(
		"refresh_token",
		response.RefreshToken,
		60*60*24*7,                   //cookie expiration time in seconds (e.g., 7 days)
		"/",                          //path where the cookie is available
		"",                           //domain
		conf.LoadConfig().Production, //secure (set to true in production to require HTTPS)
		true,                         // HttpOnly (prevents JavaScript access to the cookie)
	)

	//hiding refresh token from user for security
	c.JSON(200, &dtos.RefreshResponse{
		AccessToken: response.AccessToken,
	})
}

func (ac *AuthController) ShowClaims(c *gin.Context) {
	claims, err := ac.GetClaims(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, claims)
}

func (ac *AuthController) LogoutUser(c *gin.Context) {

	refresh, err := ac.getRefreshCookie(c)
	if err != nil {
		return
	}

	claims, err := ac.AuthService.ParseJWT(refresh)
	if err != nil {
		return
	}

	e := ac.AuthService.Logout(claims.UserID)
	if e != nil {
		c.JSON(e.Code, e.ToJson())
		return
	}

	//clear refresh token cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		conf.LoadConfig().Production,
		true,
	)

	c.Status(http.StatusNoContent)
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	refresh, err := ac.getRefreshCookie(c)
	if err != nil {
		return
	}
	response, e := ac.AuthService.RefreshToken(refresh)
	if e != nil {
		c.JSON(e.Code, e.ToJson())
		return
	}

	c.JSON(http.StatusOK, response)
}

// HELPER FUNCTIONS FOR HANDLING HEADERS/PARAMS
func (ac *AuthController) GetClaims(c *gin.Context) (*dtos.CustomClaims, error) {
	token, err := ac.ExtractAuthorization(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}

	claims, err := ac.AuthService.ParseJWT(*token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}

	return claims, nil
}

func (ac *AuthController) ExtractAuthorization(c *gin.Context) (*string, error) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return nil, fmt.Errorf("authorization header is required")
	}

	// Bearer token starts with "Bearer "
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header format")
	}

	token := parts[1]

	// Now you have the JWT token, you can process it further
	return &token, nil
}

func (ac *AuthController) ParseUserID(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return 0, err
	}

	// Ensure the ID is non-negative before casting
	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID cannot be negative",
		})
		return 0, fmt.Errorf("user ID cannot be negative")
	}

	// Cast the ID to uint
	return uint(id), nil
}

func (ac *AuthController) getRefreshCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh_token cookie"})
		return "", err
	}

	return cookie, nil
}
