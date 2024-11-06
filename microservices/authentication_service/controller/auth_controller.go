package controller

import (
	// "fmt"
	conf "authentication-service/config"
	"authentication-service/dtos"
	. "authentication-service/service"
	"fmt"
	"net/http"
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

// define all our routes for the controller
func (ac *AuthController) DefineRoutes(r *gin.Engine) {
	userGroup := r.Group("/auth")
	{
		userGroup.GET("/health", ac.TestConnection)
		userGroup.GET("/", ac.CheckIsAuthenticated)
		userGroup.POST("/login", ac.LoginUser)
		userGroup.GET("/claims", ac.ShowClaims)
		userGroup.POST("/logout", ac.LogoutUser)
		userGroup.GET("/refresh", ac.RefreshToken)
	}
}

// function meant to test connection to the service
func (ac *AuthController) TestConnection(c *gin.Context) {
	c.JSON(http.StatusOK, "This is the auth service")
}

// login user
func (ac *AuthController) LoginUser(c *gin.Context) {
	var request dtos.UserLogin

	//bind request to our dto
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	//login using auth service
	response, e := ac.AuthService.UserLogin(&request)
	if e != nil {
		c.JSON(e.Code, e.ToJson())
		return
	}

	//set our http only cookie for refresh token
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

// shows claims for a given token, assuming its in the authorization header
func (ac *AuthController) ShowClaims(c *gin.Context) {
	//get claims from token
	claims, err := ac.GetClaims(c)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, claims)
}

// endpoint used to logout user
func (ac *AuthController) LogoutUser(c *gin.Context) {
	//get refresh token from cookie
	refresh, err := ac.getRefreshCookie(c)
	if err != nil {
		return
	}

	//parse claims of token, also ensures token is valid
	claims, err := ac.AuthService.ParseJWT(&refresh)
	if err != nil {
		return
	}

	//logout user
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

// function generates new auth token from refresh token stores in cookie
func (ac *AuthController) RefreshToken(c *gin.Context) {

	refresh, err := ac.getRefreshCookie(c)
	if err != nil {
		return
	}

	//refresh auth token from cookie
	response, e := ac.AuthService.RefreshToken(refresh)
	if e != nil {
		c.JSON(e.Code, e.ToJson())
		return
	}

	c.JSON(http.StatusOK, response)
}

// endpoint to check if user is authenticated
func (ac *AuthController) CheckIsAuthenticated(c *gin.Context) {
	//extract auth token from header
	var err error
	token, err := ac.ExtractAuthorization(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	//parse claims, which checks if token is valid
	_, err = ac.AuthService.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	//able to properly extract claims which means token is valid
	c.String(http.StatusOK, "authorized")
}

// ---------- HELPER FUNCTIONS FOR HANDLING HEADERS/PARAMS -----------------------

// grabs claims from the auth token in the header, if it exists and is valid
func (ac *AuthController) GetClaims(c *gin.Context) (*dtos.CustomClaims, error) {
	token, err := ac.ExtractAuthorization(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}

	claims, err := ac.AuthService.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return nil, err
	}

	return claims, nil
}

// grabs token from authorization header if it exists
func (ac *AuthController) ExtractAuthorization(c *gin.Context) (*string, error) {
	//get the Authorization header
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return nil, fmt.Errorf("authorization header is required")
	}

	//Bearer token starts with "Bearer "
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid Authorization header format")
	}

	token := parts[1]

	//now you have the JWT token, you can process it further
	return &token, nil
}

// gets refresh token from cookie
func (ac *AuthController) getRefreshCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh_token cookie"})
		return "", err
	}

	return cookie, nil
}
