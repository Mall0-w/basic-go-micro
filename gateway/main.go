package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

// createUserServiceProxy creates a reverse proxy to the user service
func createUserServiceProxy() gin.HandlerFunc {
	// URL of the user service (matching the service name in docker-compose)
	userServiceURL, err := url.Parse("http://user-service:8080")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func CreateRouter() *gin.Engine {
	router := gin.Default()

	// Base route
	router.GET("/", helloWorld)

	// User service routes
	users := router.Group("/users")
	{
		users.Any("/*path", createUserServiceProxy())
	}

	return router
}

func main() {
	router := CreateRouter()
	// Note: Changed to 0.0.0.0 to be accessible from other containers
	router.Run("0.0.0.0:8080")
}
