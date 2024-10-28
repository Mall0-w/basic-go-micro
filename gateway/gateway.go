package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello World!")
}

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", helloWorld)
	return router
}
