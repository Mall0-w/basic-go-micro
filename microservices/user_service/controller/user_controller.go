package controller

import (
	"net/http"
	"strconv"
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
	userGroup := r.Group("/users") // Using plural form as per REST conventions
	{
		userGroup.GET("/", uc.TestConnection)
		userGroup.GET("/:id", uc.GetUserByID)
		// userGroup.POST("", uc.CreateUser)
		// userGroup.PUT("/:id", uc.UpdateUser)    // Added for completeness
		// userGroup.DELETE("/:id", uc.DeleteUser) // Added for completeness
	}
}

func (uc *UserController) TestConnection(c *gin.Context) {
	c.JSON(http.StatusOK, "This is the users service")
}

// GetUserByID handles GET requests for a single user
func (uc *UserController) GetUserByID(c *gin.Context) {
	id, err := uc.parseUserID(c)
	if err != nil {
		return
	}

	user, err := uc.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
// func (uc *UserController) CreateUser(c *gin.Context) {
// 	var request struct {
// 		Name  string `json:"name" binding:"required"`
// 		Email string `json:"email" binding:"required,email"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid request payload",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	user, err := uc.userService.CreateUser(request.Name, request.Email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to create user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, user)
// }

// // UpdateUser handles PUT requests to update an existing user
// func (uc *UserController) UpdateUser(c *gin.Context) {
// 	id, err := uc.parseUserID(c)
// 	if err != nil {
// 		return
// 	}

// 	var request struct {
// 		Name  string `json:"name" binding:"required"`
// 		Email string `json:"email" binding:"required,email"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "Invalid request payload",
// 			"details": err.Error(),
// 		})
// 		return
// 	}

// 	user, err := uc.userService.UpdateUser(id, request.Name, request.Email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to update user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// // DeleteUser handles DELETE requests to remove a user
// func (uc *UserController) DeleteUser(c *gin.Context) {
// 	id, err := uc.parseUserID(c)
// 	if err != nil {
// 		return
// 	}

// 	if err := uc.userService.DeleteUser(id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to delete user",
// 		})
// 		return
// 	}

// 	c.Status(http.StatusNoContent)
// }

// // parseUserID is a helper function to parse and validate user IDs from requests
func (uc *UserController) parseUserID(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return 0, err
	}
	return id, nil
}
