package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// simple function to test connection to gateway
func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

// helper function to create a proxy for a url
func CreateProxy(proxyUrl string) gin.HandlerFunc {
	// URL of the user service (matching the service name in docker-compose)
	serviceURL, err := url.Parse(proxyUrl)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(serviceURL)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// abstracting create router for testing
func CreateRouter() *gin.Engine {
	router := gin.Default()

	//Base Router
	router.GET("/", helloWorld)

	//Creating groups and proxing them to different services based on path
	users := router.Group("/users")
	{
		users.Any("/*path", CreateProxy("http://user-service:8080"))
	}

	inventory := router.Group("/inventory")
	{
		inventory.Any("/*path", CreateProxy("http://inventory-service:8080"))
	}

	auth := router.Group("/auth")
	{
		auth.Any("/*path", CreateProxy("http://auth-service:8080"))
	}

	return router
}

func main() {
	router := CreateRouter()
	//note: Changed to 0.0.0.0 to be accessible from other containers
	router.Run("0.0.0.0:8080")
}
