package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/dtos"
	. "user-service/service"

	"github.com/gin-gonic/gin"
)

// UserController handles user-related HTTP requests
type UserController struct {
	userService *UserService // Assuming you have a UserService
}

// NewUserController creates a new UserController instance
func NewUserController(userService *UserService) *UserController {
	if userService == nil {
		userService = NewUserService(nil)
	}

	return &UserController{
		userService: userService,
	}
}

// DefineRoutes sets up the routing for user endpoints
func (uc *UserController) DefineRoutes(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.GET("/health", uc.TestConnection)
		userGroup.GET("/:id", uc.GetUserByID)
		userGroup.POST("/", uc.CreateUser)
		userGroup.PUT("/:id", uc.UpdateUser)
		userGroup.DELETE("/:id", uc.DeleteUser)
	}
}

func (uc *UserController) TestConnection(c *gin.Context) {
	c.JSON(http.StatusOK, "This is the users service")
}

// GetUserByID handles GET requests for a single user
func (uc *UserController) GetUserByID(c *gin.Context) {
	id, e := uc.parseUserID(c)
	if e != nil {
		return
	}

	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.ToJson())
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
func (uc *UserController) CreateUser(c *gin.Context) {

	var request dtos.UserCreate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	user, e := uc.userService.CreateUser(request)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e.ToJson())
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser handles PUT requests to update an existing user
func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := uc.parseUserID(c)
	if err != nil {
		return
	}

	var request dtos.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}

	request.Id = id

	user, e := uc.userService.UpdateUser(request)
	if e != nil {
		c.JSON(http.StatusInternalServerError, e.ToJson())
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE requests to remove a user
func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := uc.parseUserID(c)
	if err != nil {
		return
	}

	if e := uc.userService.DeleteUser(id); e != nil {
		c.JSON(http.StatusInternalServerError, e.ToJson())
		return
	}

	c.Status(http.StatusNoContent)
}

// parseUserID is a helper function to parse and validate user IDs from requests
func (uc *UserController) parseUserID(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return 0, err
	}

	//ensure the ID is non-negative before casting
	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID cannot be negative",
		})
		return 0, fmt.Errorf("user ID cannot be negative")
	}

	//cast the ID to uint
	return uint(id), nil
}
